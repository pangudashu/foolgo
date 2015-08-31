package foolgo

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	STATE_INIT          = 1
	STATE_RUNNING       = 2
	STATE_TERMINATE     = 3
	STATE_SHUTTING_DOWN = 4
)

var (
	runLock    *sync.Mutex
	connWg     sync.WaitGroup
	isChild    bool
	operate    string
	serverStat int
	isForking  bool
	restart    string
	logger     *Log
)

type HttpServerConfig struct {
	RunMod        string
	Root          string //web访问目录
	ViewPath      string
	Addr          string
	AccessLog     string
	ErrorLog      string
	StandOutPut   string
	IsGzip        bool
	ZipMinSize    int
	ReadTimeout   int
	WriteTimeout  int
	MaxHeaderByte int
	HttpErrorHtml map[int]string
	//https
	SslOn      bool
	SslCert    string
	SslCertKey string
	Pid        string
}

type FoolServer struct {
	*http.Server
	listener    net.Listener
	listenerPtr *FoolListener
	App         *Application
	config      *HttpServerConfig
}

func init() {
	runLock = &sync.Mutex{}

	flag.BoolVar(&isChild, "reload", false, "listen on open fd (after forking)")

	var cmd *string = flag.String("s", "start", "[cmd:restart|stop]")
	flag.Parse()
	operate = *cmd

	//防止重复fork
	isForking = false

	if operate == "restart" {
	} else if operate == "stop" {
	}
}

func NewServer(server_config *HttpServerConfig) (*FoolServer, error) {
	runLock.Lock()
	defer runLock.Unlock()

	if server_config.Addr == "" {
		return nil, errors.New("server Addr can't be empty...[ip:port]")
	}
	if server_config.ReadTimeout <= 0 {
		server_config.ReadTimeout = 30
	}
	if server_config.WriteTimeout <= 0 {
		server_config.WriteTimeout = 30
	}
	if server_config.MaxHeaderByte <= 0 {
		server_config.MaxHeaderByte = 1 << 20
	}
	if server_config.Pid == "" {
		return nil, errors.New("foolgo.HttpServerConfig.Pid can't be empty")
	}

	srv := &FoolServer{config: server_config}

	l, err := NewListener(server_config.Addr)
	if err != nil {
		return nil, err
	}

	if server_config.SslOn == true && server_config.SslCert != "" && server_config.SslCertKey != "" {
		srv.listenerPtr = l
		srv.listener, err = NewTlsListener(srv.listenerPtr, server_config.SslCert, server_config.SslCertKey)
		if err != nil {
			return nil, err
		}
	} else {
		srv.listenerPtr = l
		srv.listener = l
	}
	//new Application
	app, err := NewApplication(server_config)
	if err != nil {
		return nil, err
	}

	srv.App = app
	srv.Server = &http.Server{}
	srv.Server.Addr = server_config.Addr
	srv.Server.ReadTimeout = time.Duration(server_config.ReadTimeout) * time.Second
	srv.Server.WriteTimeout = time.Duration(server_config.WriteTimeout) * time.Second
	srv.Server.MaxHeaderBytes = server_config.MaxHeaderByte
	srv.Server.Handler = app

	return srv, nil
}

func (srv *FoolServer) RegRewrite(rewrite map[string]string) *FoolServer {
	regRewrite(rewrite)
	return srv
}

func (srv *FoolServer) Run() {
	logger = NewLog(srv.config.AccessLog, srv.config.ErrorLog, srv.config.StandOutPut)
	//解析模板
	CompileTpl(srv.config.ViewPath)
	//信号处理函数
	go srv.signalHandle()

	serverStat = STATE_RUNNING

	//kill父进程
	if isChild == true {
		parent := syscall.Getppid()

		if _, err := os.FindProcess(parent); err != nil {
			return
		}
		log.Printf("main: Killing parent pid: %v", parent)
		syscall.Kill(parent, syscall.SIGQUIT)
	}
	srv.createPid(syscall.Getpid())

	//listen loop
	srv.Serve(srv.listener)

	log.Println(syscall.Getpid(), "[server.go]Waiting for connections to finish...")
	connWg.Wait()
	serverStat = STATE_TERMINATE
	log.Println("[server.go]server shuttdown!!!!")
	return
}

//信号处理
func (srv *FoolServer) signalHandle() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		sig := <-ch

		switch sig {
		case syscall.SIGHUP:
			srv.forkServer()
		case syscall.SIGINT:
			srv.shutDown()
		case syscall.SIGQUIT:
			srv.shutDown()
		case syscall.SIGTERM:
			srv.shutDown()
		default:
		}
	}
}

func (srv *FoolServer) shutDown() {
	if serverStat != STATE_RUNNING {
		return
	}
	serverStat = STATE_SHUTTING_DOWN

	go srv.serverTimeout()

	err := srv.listener.Close()
	if err != nil {
		log.Println(syscall.Getpid(), "Listener.Close() error:", err)
	} else {
		log.Println(syscall.Getpid(), "[server.go#shutDown]", srv.listener.Addr(), "Listener closed.")
	}
}

func (srv *FoolServer) serverTimeout() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("WaitGroup at 0", r)
		}
	}()
	if serverStat != STATE_SHUTTING_DOWN {
		return
	}
	time.Sleep(time.Second * 60)
	for {
		if serverStat == STATE_TERMINATE {
			break
		}
		connWg.Done()
	}
	log.Println("[STOP - Hammer Time] Forcefully shutting down parent")
}

func (srv *FoolServer) createPid(pid int) {
	fd, _ := os.OpenFile(srv.config.Pid, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer fd.Close()
	fd.WriteString(fmt.Sprintf("%d", pid))

	_, err := os.Stat(srv.config.Pid)
	if err != nil && os.IsNotExist(err) {
		fmt.Printf("Can't create %q\n", srv.config.Pid)
		os.Exit(1)
	}
}

//重启server
func (srv *FoolServer) forkServer() {
	runLock.Lock()
	defer runLock.Unlock()

	if isForking {
		return
	}
	isForking = true

	var file *os.File
	file = srv.listenerPtr.File()

	files := make([]*os.File, 1)
	files[0] = file

	path := os.Args[0]
	var args []string
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			if arg == "-reload" {
				continue
			}
			args = append(args, arg)
		}
	}
	args = append(args, "-reload")

	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = files

	err := cmd.Start()
	if err != nil {
		log.Fatalf("Restart: Failed to launch, error: %v", err)
	}

	return
}

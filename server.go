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
)

type HttpServerConfig struct {
	RunMod        string
	Root          string //web访问目录
	ViewPath      string
	Addr          string
	IsGzip        bool
	ZipMinSize    int
	ReadTimeout   int
	WriteTimeout  int
	MaxHeaderByte int
	HttpErrorHtml map[int]string
}

type FoolServer struct {
	*http.Server
	listener net.Listener
	App      *Application
	config   *HttpServerConfig
}

var restart string

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

	l, err := NewListener(server_config.Addr)
	if err != nil {
		return nil, err
	}
	//new Application
	app, err := NewApplication(server_config)
	if err != nil {
		return nil, err
	}

	srv := &FoolServer{listener: l, App: app, config: server_config}
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
	//解析模板
	CompileTpl(srv.config.ViewPath)

	//信号处理函数
	go srv.signalHandle()

	serverStat = STATE_RUNNING

	//杀掉父进程
	if isChild == true {
		parent := syscall.Getppid()

		if _, err := os.FindProcess(parent); err != nil {
			return
		}
		log.Printf("main: Killing parent pid: %v", parent)
		syscall.Kill(parent, syscall.SIGTERM)
	}

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
			log.Print("signal hup")
			srv.forkServer()
		case syscall.SIGINT:
			srv.shutDown()
		case syscall.SIGQUIT:
			srv.shutDown()
		case syscall.SIGTERM:
			srv.shutDown()
		default:
			log.Print("unknown")
		}
	}
}

//关闭server
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
	time.Sleep(time.Second * 20)
	log.Println("[STOP - Hammer Time] Forcefully shutting down parent")
	for {
		if serverStat == STATE_TERMINATE {
			break
		}
		connWg.Done()
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

	file := srv.listener.(*FoolListener).File()
	fmt.Println(file)

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

package foolgo

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	application_instance   *Application = nil
	DEFAULT_CONTROLLER     string       = "index"
	DEFAULT_ACTION         string       = "index"
	ACTION_SUFFIX          string       = "Action"
	HTTP_METHOD_PARAM_NAME string       = "m"
	IsGzip                 bool         = true
	ZipMinSize             int          = 1024
	RunMod                 string       = "product"
)

type Application struct {
	server_config *HttpServerConfig
	register      *Register
	dispatcher    *Dispatcher
}

//create Application object
func NewApplication(http_server_config *HttpServerConfig) (*Application, error) { /*{{{*/
	if application_instance != nil {
		return application_instance, nil
	}

	if http_server_config == nil {
		return nil, fmt.Errorf("please init Http_Server_Config")
	}

	http_server_config.Root = strings.TrimRight(http_server_config.Root, "/")
	for err_code, err_file_name := range http_server_config.HttpErrorHtml {
		err_html := http_server_config.Root + "/" + strings.TrimLeft(err_file_name, "/")
		http_server_config.HttpErrorHtml[err_code] = err_html
	}

	application_instance = &Application{
		server_config: http_server_config,
	}
	RunMod = http_server_config.RunMod
	IsGzip = http_server_config.IsGzip
	ZipMinSize = http_server_config.ZipMinSize

	//init mime
	initMime()
	//init Register
	application_instance.register = NewRegister()
	//init dispatcher
	application_instance.dispatcher = NewDispatcher(http_server_config)
	//init router
	NewRouter()

	return application_instance, nil
} /*}}}*/

func (this *Application) AddViewFunc(key string, func_name interface{}) {
	AddViewFunc(key, func_name)
}

func (this *Application) AddCompressType(file_ext []string) {
	for _, v := range file_ext {
		AddCompressType(v)
	}
}

//register controller
func (this *Application) RegController(register_controller_map map[string]FGController) { /*{{{*/
	for controller_name, controller := range register_controller_map {
		err := this.register.SetController(controller_name, controller)
		if err != nil {
			logger.RunLog(fmt.Sprintf("[Error] RegController error :%v", err))
			os.Exit(0)
		}
	}
} /*}}}*/

//http request handler
func (this *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) { /*{{{*/
	start_time := time.Now()

	if r.URL.Path != "/" {
		matchRewrite(r)
	}
	this.dispatcher.Dispatch_handler(w, r)

	end_time := time.Now()

	request_time := float64(end_time.UnixNano()-start_time.UnixNano()) / 1000000000

	log_format := "%s - [%s] %s %s %s %s %.5f \"%s\"" //ip - [time] method uri scheme status request_time agent
	access_log := fmt.Sprintf(log_format,
		this.Isset(r.RemoteAddr),
		Date("Y/m/d H:i:s", start_time),
		this.Isset(r.Method),
		this.Isset(r.URL.RequestURI()),
		this.Isset(r.Proto),
		this.Isset(w.Header().Get("Status")),
		request_time,
		this.Isset(r.Header.Get("User-Agent")),
	)
	logger.AccessLog(access_log)
} /*}}}*/

func (this *Application) Isset(params string) string { /*{{{*/
	if params == "" {
		return "-"
	}
	return params
} /*}}}*/

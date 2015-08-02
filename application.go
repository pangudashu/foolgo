package foolgo

import (
	"fmt"
	"net/http"
	"os"
	"strings"
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

/*{{{ func NewApplication(http_server_config *HttpServerConfig) (*Application,error)
 *获取Application对象
 */
func NewApplication(http_server_config *HttpServerConfig) (*Application, error) {
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

	//初始化mime
	initMime()
	//初始化Register
	application_instance.register = NewRegister()
	//初始化dispatcher
	application_instance.dispatcher = NewDispatcher(http_server_config)
	//初始化router
	NewRouter()
	//解析模板
	CompileTpl(http_server_config.ViewPath)

	return application_instance, nil
}

/*}}}*/

/*{{{ func (this *Application) RegController(controller_name string, controller interface{})
 */
func (this *Application) RegController(register_controller_map map[string]FGController) {
	for controller_name, controller := range register_controller_map {
		err := this.register.SetController(controller_name, controller)
		if err != nil {
			fmt.Println("[RegController]:", err)
			os.Exit(0)
		}
	}
}

/*}}}*/

/*{{{ func (this *Application) ServeHTTP(w http.ResponseWriter, r *http.Request)
 * Http请求入口
 */
func (this *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		matchRewrite(r)
	}
	this.dispatcher.Dispatch_handler(w, r)
}

/*}}}*/

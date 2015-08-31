package foolgo

import (
	//"fmt"
	"mime/multipart"
	"net/http"
)

type Controller struct {
	request  *Request
	response *Response
	view     *View
}

type FGController interface {
	RegRouter() map[string]interface{}
}

func (this *Controller) Init(request *Request, response *Response) bool {
	this.request = request
	this.response = response
	this.view = NewView()

	return true
}

//no use just implement FGController
func (this *Controller) RegRouter() map[string]interface{} {
	return nil
}

/*{{{ func (this *Controller) Param(key string, default_value ...string) string
 */
func (this *Controller) Param(key string, default_value ...string) string {
	v := this.request.Param(key)
	if v == "" && default_value != nil {
		return default_value[0]
	}
	return v
}

/*}}}*/

/*{{{ func (this *Controller) Assign(key interface{}, value interface{})
 */
func (this *Controller) Assign(key interface{}, value interface{}) {
	this.view.Assign(key, value)
}

/*}}}*/

/*{{{ func (this *Controller) Display(view_path ...string)
 */
func (this *Controller) Display(view_path ...string) {
	bytes, err := this.Render(view_path...)

	if err == nil {
		this.response.Header("Content-Type", "text/html; charset=utf-8")
		this.response.Body(bytes)
	} else {
		this.response.Header("Content-Type", "text/html; charset=utf-8")
		this.response.Body([]byte(err.Error()))
	}
}

/*}}}*/

/*{{{ func (this *Controller) Render(view_path ...string) ([]byte, error)
 */
func (this *Controller) Render(view_path ...string) ([]byte, error) {
	var view_name string
	if view_path == nil || view_path[0] == "" {
		view_name = this.request.GetController() + "/" + this.request.GetAction()
	} else {
		view_name = view_path[0]
	}
	return this.view.Render(view_name)
}

/*}}}*/

func (this *Controller) Cookie(name string) string {
	return this.request.Cookie(name)
}

func (this *Controller) Uri() string {
	return this.request.Uri()
}

func (this *Controller) Url() string {
	return this.request.Url()
}

func (this *Controller) IP() string {
	return this.request.IP()
}

func (this *Controller) Scheme() string {
	return this.request.Scheme()
}

func (this *Controller) Header(key string) string {
	return this.request.Header(key)
}

func (this *Controller) SetHeader(key, value string) {
	this.response.Header(key, value)
}

func (this *Controller) SetCookie(name string, value string, others ...interface{}) {
	this.response.Cookie(name, value, others...)
}

func (this *Controller) OutString(bytes []byte) {
	this.response.Header("Content-Type", "text/html; charset=utf-8")
	this.response.Body(bytes)
}

func (this *Controller) Json(data interface{}, coding ...bool) error {
	return this.response.Json(data, coding...)
}

func (this *Controller) Jsonp(callback string, data interface{}, coding ...bool) error {
	return this.response.Jsonp(callback, data, coding...)
}

func (this *Controller) Method() string {
	return this.request.Method()
}

//获取所有get变量
func (this *Controller) GET() map[string]string {
	return this.request.ParamGet()
}

//获取所有post提交变量
func (this *Controller) POST() map[string]interface{} {
	return this.request.ParamPost()
}

func (this *Controller) Location(url string) {
	http.Redirect(this.response.Writer, this.request.request, url, 301)
}

// 获取所有上传文件
// files, _ := this.GetUploadFiles("user_icon")
// for i, _ := range files {
//	 file, _ := files[i].Open()
//	 defer file.Close()
//	 log.Print(this.GetFileSize(&file))
// }
func (this *Controller) GetUploadFiles(key string) ([]*multipart.FileHeader, error) {
	return this.request.GetUploadFiles(key)
}

func (this *Controller) MoveUploadFile(fromfile, tofile string) error {
	return this.request.MoveUploadFile(fromfile, tofile)
}

func (this *Controller) GetFileSize(file *multipart.File) int64 {
	return this.request.GetFileSize(file)
}

package foolgo

import (
	"errors"
	//"fmt"
	"mime/multipart"
	"net/http"
	//"reflect"
	"io"
	"os"
	"strings"
)

type Request struct {
	request         *http.Request
	controller_name string
	action_name     string
	form_parsed     bool
	rewrite_params  map[string]string
}

func NewRequest(r *http.Request) *Request {
	request_instance := &Request{
		request: r,
	}

	return request_instance
}

func (this *Request) SetController(name string) {
	if name == "" {
		return
	}
	this.controller_name = name
}

func (this *Request) SetAction(name string) {
	if name == "" {
		return
	}
	this.action_name = name
}

func (this *Request) GetController() string {
	return this.controller_name
}

func (this *Request) GetAction() string {
	return this.action_name
}

func (this *Request) Method() string {
	return this.request.Method
}

func (this *Request) Uri() string {
	return this.request.RequestURI
}

func (this *Request) Url() string {
	return this.request.URL.Path
}

func (this *Request) Header(key string) string {
	return this.request.Header.Get(key)
}

func (this *Request) Cookie(key string) string {
	cookie, err := this.request.Cookie(key)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func (this *Request) Param(key string) string {
	if this.rewrite_params != nil {
		if d, ok := this.rewrite_params[key]; ok == true {
			return d
		}
	}
	err := this.ParseMultiForm()
	if err != nil {
		return ""
	}
	return this.request.Form.Get(key)
}

func (this *Request) ParamGet() (data map[string]string) {
	err := this.ParseMultiForm()
	if err != nil {
		return nil
	}

	if this.request.Form == nil {
		return this.rewrite_params
	}
	data = make(map[string]string)
	for k, v := range this.request.Form {
		data[k] = v[0]
	}
	if this.rewrite_params != nil {
		for k, v := range this.rewrite_params {
			data[k] = v
		}
	}
	return data
}

func (this *Request) ParamPost() (data map[string]interface{}) {
	err := this.ParseMultiForm()
	if err != nil {
		return nil
	}
	if this.request.PostForm == nil {
		return nil
	}
	data = make(map[string]interface{})
	for k, v := range this.request.PostForm {
		if len(v) > 1 {
			data[k] = v
		} else {
			data[k] = v[0]
		}
	}
	return data
}

/*{{{ func (this *Request) ParseForm() error
 */
func (this *Request) ParseMultiForm() error {
	if this.form_parsed == true || this.request.Form != nil || this.request.PostForm != nil || this.request.MultipartForm != nil {
		return nil
	}

	this.form_parsed = true
	if strings.Contains(this.Header("Content-Type"), "multipart/form-data") {
		if err := this.request.ParseMultipartForm(32 << 20); err != nil {
			return errors.New("Error parsing request body:" + err.Error())
		}
	} else if err := this.request.ParseForm(); err != nil {
		return errors.New("Error parsing request body:" + err.Error())
	}
	return nil
}

/*}}}*/

func (this *Request) GetUploadFiles(key string) ([]*multipart.FileHeader, error) {
	this.ParseMultiForm()

	if this.request.MultipartForm == nil {
		return nil, nil
	}
	files, ok := this.request.MultipartForm.File[key]
	if ok {
		return files, nil
	}
	return nil, http.ErrMissingFile
}

func (this *Request) MoveUploadFile(fromfile, tofile string) error {
	file, _, err := this.request.FormFile(fromfile)
	if err != nil {
		return err
	}

	defer file.Close()

	f, err := os.OpenFile(tofile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	io.Copy(f, file)
	return nil
}

type Size interface {
	Size() int64
}

func (this *Request) GetFileSize(file *multipart.File) int64 {
	if sizeInterface, ok := (*file).(Size); ok {
		return sizeInterface.Size()
	}
	return -1
}

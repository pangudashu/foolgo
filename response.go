package foolgo

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Response struct {
	Writer        http.ResponseWriter
	request       *Request
	server_config *HttpServerConfig
}

func NewResponse(w http.ResponseWriter, request *Request, config *HttpServerConfig) *Response {
	response := &Response{
		Writer:        w,
		request:       request,
		server_config: config,
	}
	return response
}

var cookieNameFilter = strings.NewReplacer("\n", "-", "\r", "-")
var cookieValueFilter = strings.NewReplacer("\n", " ", "\r", " ", ";", " ")

func filterName(n string) string {
	return cookieNameFilter.Replace(n)
}

func filterValue(v string) string {
	return cookieValueFilter.Replace(v)
}

// Set http response header
func (this *Response) Header(key, val string) {
	this.Writer.Header().Set(key, val)
}

func (this *Response) Body(html_content []byte) { /*{{{*/
	accept_encoding := this.request.Header("Accept-Encoding")
	if IsGzip == true && len(html_content) >= ZipMinSize && accept_encoding != "" && strings.Index(accept_encoding, "gzip") >= 0 {
		this.Header("Content-Encoding", "gzip")

		output_writer, _ := gzip.NewWriterLevel(this.Writer, gzip.BestSpeed)
		defer output_writer.Close()
		output_writer.Write(html_content)
	} else {
		this.Writer.Write(html_content)
	}
} /*}}}*/

// Set cookie
// Copy from beego @https://github.com/astaxie/beego
func (this *Response) Cookie(name string, value string, others ...interface{}) { /*{{{*/
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s=%s", filterName(name), filterValue(value))
	if len(others) > 0 {
		switch v := others[0].(type) {
		case int:
			if v > 0 {
				fmt.Fprintf(&b, "; Max-Age=%d", v)
			} else if v < 0 {
				fmt.Fprintf(&b, "; Max-Age=0")
			}
		case int64:
			if v > 0 {
				fmt.Fprintf(&b, "; Max-Age=%d", v)
			} else if v < 0 {
				fmt.Fprintf(&b, "; Max-Age=0")
			}
		case int32:
			if v > 0 {
				fmt.Fprintf(&b, "; Max-Age=%d", v)
			} else if v < 0 {
				fmt.Fprintf(&b, "; Max-Age=0")
			}
		}
	}
	// the settings below
	// Path, Domain, Secure, HttpOnly
	// can use nil skip set

	// default "/"
	if len(others) > 1 {
		if v, ok := others[1].(string); ok && len(v) > 0 {
			fmt.Fprintf(&b, "; Path=%s", filterValue(v))
		}
	} else {
		fmt.Fprintf(&b, "; Path=%s", "/")
	}

	if len(others) > 2 {
		if v, ok := others[2].(string); ok && len(v) > 0 {
			fmt.Fprintf(&b, "; Domain=%s", filterValue(v))
		}
	}

	if len(others) > 3 {
		var secure bool
		switch v := others[3].(type) {
		case bool:
			secure = v
		default:
			if others[3] != nil {
				secure = true
			}
		}
		if secure {
			fmt.Fprintf(&b, "; Secure")
		}
	}

	// default false. for session cookie default true
	httponly := false
	if len(others) > 4 {
		if v, ok := others[4].(bool); ok && v {
			// HttpOnly = true
			httponly = true
		}
	}

	if httponly {
		fmt.Fprintf(&b, "; HttpOnly")
	}

	this.Writer.Header().Add("Set-Cookie", b.String())
} /*}}}*/

// Set output type:json
func (this *Response) Json(data interface{}, coding ...bool) error { /*{{{*/
	this.Header("Content-Type", "application/json;charset=UTF-8")

	var content []byte
	var err error

	content, err = json.Marshal(data)
	if err != nil {
		http.Error(this.Writer, err.Error(), http.StatusInternalServerError)
		return err
	}
	if coding != nil && coding[0] == true {
		content = []byte(unicode(string(content)))
	}
	this.Body(content)
	return nil
} /*}}}*/

func (this *Response) Jsonp(callback string, data interface{}, coding ...bool) error { /*{{{*/
	this.Header("Content-Type", "application/javascript;charset=UTF-8")

	var content []byte
	var err error

	content, err = json.Marshal(data)
	if err != nil {
		http.Error(this.Writer, err.Error(), http.StatusInternalServerError)
		return err
	}
	if coding != nil && coding[0] == true {
		content = []byte(unicode(string(content)))
	}
	ck := bytes.NewBufferString(" " + callback)
	ck.WriteString("(")
	ck.Write(content)
	ck.WriteString(");\r\n")

	this.Body(ck.Bytes())
	return nil
} /*}}}*/

// Convert to unicode
func unicode(str string) string { /*{{{*/
	rs := []rune(str)
	jsons := ""
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			jsons += string(r)
		} else {
			jsons += "\\u" + strconv.FormatInt(int64(rint), 16) // json
		}
	}
	return jsons
} /*}}}*/

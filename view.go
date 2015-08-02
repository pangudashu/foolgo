package foolgo

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	ViewRoot       string
	ViewExt        string = ".html"
	ViewTemplates  map[string]*template.Template
	template_files map[string]string
)

type View struct {
	data map[interface{}]interface{}
}

func NewView() *View {
	return &View{}
}

/*{{{ func (this *View) Assign(key interface{}, value interface{})
 */
func (this *View) Assign(key interface{}, value interface{}) {
	if this.data == nil {
		this.data = make(map[interface{}]interface{})
		this.data[key] = value
	} else {
		this.data[key] = value
	}
}

/*}}}*/

/*{{{ func (this *View) Render(view_name string) (string, error)
 */
func (this *View) Render(view_name string) ([]byte, error) {
	view_name = strings.ToLower(view_name)
	if RunMod == "dev" {
		t := template.New(view_name).Delims("{{", "}}")
		t, err := parseTemplate(t, ViewRoot+"/"+view_name+ViewExt)
		if err != nil || t == nil {
			return []byte(""), err
		}
		ViewTemplates[view_name] = t
	}

	if tpl, ok := ViewTemplates[view_name]; ok == false {
		return []byte(""), errors.New("template " + view_name + " not complile")
	} else {
		html_content_bytes := bytes.NewBufferString("")
		err := tpl.ExecuteTemplate(html_content_bytes, view_name, this.data)
		if err != nil {
			return []byte(""), err
		}
		html_content, _ := ioutil.ReadAll(html_content_bytes)
		return html_content, nil

		/*
			if IsGzip == false {
				return string(html_content), nil
			}

			var b bytes.Buffer
			w := gzip.NewWriter(&b)
			defer w.Close()

			w.Write(html_content)
			w.Flush()

			htmls, _ := ioutil.ReadAll(&b)
			return string(htmls), nil
		*/
	}
}

/*}}}*/

/*{{{ func CompileTpl(view_root string) error
 */
func CompileTpl(view_root string) error {
	ViewRoot = view_root
	template_files = make(map[string]string)

	filepath.Walk(view_root, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() || (f.Mode()&os.ModeSymlink) > 0 {
			return nil
		}

		if strings.HasSuffix(path, ViewExt) == false {
			return nil
		}

		file_name := strings.Trim(strings.Replace(path, ViewRoot, "", 1), "/")
		template_files[strings.TrimSuffix(file_name, ViewExt)] = path
		return nil
	})

	ViewTemplates = make(map[string]*template.Template)

	for name, file := range template_files {
		if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
			fmt.Printf("parse template %q err : %q", file, err)
			continue
		}
		t := template.New(name).Delims("{{", "}}")

		t, err := parseTemplate(t, file)
		if err != nil || t == nil {
			continue
		}
		ViewTemplates[name] = t
	}

	return nil
}

/*}}}*/

/*{{{ func parseTemplate(template *template.Template, file string) (t *template.Template, err error)
 */
func parseTemplate(template *template.Template, file string) (t *template.Template, err error) {
	data, _ := ioutil.ReadFile(file)

	t, err = template.Parse(string(data))
	if err != nil {
		return nil, err
	}

	reg := regexp.MustCompile(`{{\s{0,}template\s{0,}"(.*?)".*?}}`)
	match := reg.FindAllStringSubmatch(string(data), -1)
	for _, v := range match {
		if v == nil || v[1] == "" {
			continue
		}
		tlook := t.Lookup(v[1])
		if tlook != nil {
			continue
		}
		deep_file := ViewRoot + "/" + v[1] + ViewExt
		if deep_file == file {
			continue
		}

		t, err = parseTemplate(t, deep_file)

		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

/*}}}*/

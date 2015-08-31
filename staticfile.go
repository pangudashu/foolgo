package foolgo

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	GzipExt = []string{".css", ".js", ".html"}
)

func OutStaticFile(response *Response, request *Request, file string) { /*{{{*/
	file_path := response.server_config.Root + file
	fi, err := os.Stat(file_path)
	if err != nil && os.IsNotExist(err) {
		OutErrorHtml(response, request, http.StatusNotFound)
		return
	} else if fi.IsDir() == true {
		OutErrorHtml(response, request, http.StatusForbidden)
		return
	}
	file_size := fi.Size()
	mod_time := fi.ModTime()

	if IsGzip == false || file_size < int64(ZipMinSize) {
		http.ServeFile(response.Writer, request.request, file_path)
		return
	}

	is_gzip := false
	for _, ext := range GzipExt {
		if strings.HasSuffix(strings.ToLower(file), strings.ToLower(ext)) {
			is_gzip = true
			break
		}
	}
	if is_gzip == false {
		http.ServeFile(response.Writer, request.request, file_path)
		return
	}
	osfile, err := os.Open(file_path)
	if err != nil {
		OutErrorHtml(response, request, http.StatusNotFound)
		return
	}

	var b bytes.Buffer
	output_writer, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
	if err != nil {
		OutErrorHtml(response, request, http.StatusNotFound)
		return
	}
	_, err = io.Copy(output_writer, osfile)
	output_writer.Close()
	if err != nil {
		OutErrorHtml(response, request, http.StatusNotFound)
		return
	}
	content, err := ioutil.ReadAll(&b)
	if err != nil {
		OutErrorHtml(response, request, http.StatusNotFound)
		return
	}
	cfi := &memFileInfo{fi, mod_time, content, int64(len(content)), file_size}
	mf := &memFile{cfi, 0}

	response.Header("Content-Encoding", "gzip")
	http.ServeContent(response.Writer, request.request, file_path, mod_time, mf)
} /*}}}*/

func OutErrorHtml(response *Response, request *Request, http_code int) { /*{{{*/
	response.Header("Status", fmt.Sprintf("%d", http_code))
	if err_html, ok := response.server_config.HttpErrorHtml[http_code]; ok == true {
		if fi, err := os.Stat(err_html); (err == nil || os.IsExist(err)) && fi.IsDir() != true {
			http.ServeFile(response.Writer, request.request, err_html)
			return
		}
	}

	error_list := make(map[int]string)

	error_list[403] = `
	<html>
	<head><title>403 Forbidden</title></head>
	<body bgcolor="white">
	<center><h1>403 Forbidden</h1></center>
	<hr><center>foolgo/1.0.0</center>
	</body>
	</html>
	<!-- a padding to disable MSIE and Chrome friendly error page -->
	<!-- a padding to disable MSIE and Chrome friendly error page -->
	<!-- a padding to disable MSIE and Chrome friendly error page -->
	<!-- a padding to disable MSIE and Chrome friendly error page -->
	<!-- a padding to disable MSIE and Chrome friendly error page -->
	<!-- a padding to disable MSIE and Chrome friendly error page -->
	`
	error_list[404] = `
	<html>
	<head><title>404 Not Found</title></head>
	<body bgcolor="white">
	<center><h1>404 Not Found</h1></center>
	<hr><center>foolgo/1.0.0</center>
	</body>
	</html>
	<!-- a padding to disable MSIE and Chrome friendly error page -->
	<!-- a padding to disable MSIE and Chrome friendly error page -->
	<!-- a padding to disable MSIE and Chrome friendly error page -->
	<!-- a padding to disable MSIE and Chrome friendly error page -->
	<!-- a padding to disable MSIE and Chrome friendly error page -->
	<!-- a padding to disable MSIE and Chrome friendly error page -->
	`
	error_list[500] = `
	<html>
	<head><title>500 Internal Server Error</title></head>
	<body bgcolor="white">
	<center><h1>500 Internal Server Error</h1></center>
	<hr><center>foolgo/1.0.0</center>
	</body>
	</html>
	<!-- a padding to disable MSIE and Chrome friendly error page -->
	<!-- a padding to disable MSIE and Chrome friendly error page -->
	<!-- a padding to disable MSIE and Chrome friendly error page -->
	<!-- a padding to disable MSIE and Chrome friendly error page -->
	<!-- a padding to disable MSIE and Chrome friendly error page -->
	<!-- a padding to disable MSIE and Chrome friendly error page -->
	`
	response.Header("Content-Type", "text/html; charset=utf-8")
	response.Header("X-Content-Type-Options", "nosniff")
	response.Writer.WriteHeader(http_code)
	fmt.Fprintln(response.Writer, error_list[http_code])
} /*}}}*/

type memFile struct {
	fi     *memFileInfo
	offset int64
}

// Close memfile.
func (f *memFile) Close() error {
	return nil
}

// Get os.FileInfo of memfile.
func (f *memFile) Stat() (os.FileInfo, error) {
	return f.fi, nil
}

// read os.FileInfo of files in directory of memfile.
// it returns empty slice.
func (f *memFile) Readdir(count int) ([]os.FileInfo, error) {
	infos := []os.FileInfo{}

	return infos, nil
}

func (f *memFile) Read(p []byte) (n int, err error) {
	if len(f.fi.content)-int(f.offset) >= len(p) {
		n = len(p)
	} else {
		n = len(f.fi.content) - int(f.offset)
		err = io.EOF
	}
	copy(p, f.fi.content[f.offset:f.offset+int64(n)])
	f.offset += int64(n)
	return
}

var errWhence = errors.New("Seek: invalid whence")
var errOffset = errors.New("Seek: invalid offset")

func (f *memFile) Seek(offset int64, whence int) (ret int64, err error) {
	switch whence {
	default:
		return 0, errWhence
	case os.SEEK_SET:
	case os.SEEK_CUR:
		offset += f.offset
	case os.SEEK_END:
		offset += int64(len(f.fi.content))
	}
	if offset < 0 || int(offset) > len(f.fi.content) {
		return 0, errOffset
	}
	f.offset = offset
	return f.offset, nil
}

type memFileInfo struct {
	os.FileInfo
	modTime     time.Time
	content     []byte
	contentSize int64
	fileSize    int64
}

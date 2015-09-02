package main

import (
	"fmt"
	"github.com/pangudashu/foolgo"
	"github.com/pangudashu/foolgo/example/application/controllers"
)

var controller_map = map[string]foolgo.FGController{
	"demo": &controllers.DemoController{},
}

func main() {
	//server config
	server_config := &foolgo.HttpServerConfig{
		Root:       "/home/qinpeng/mygo/src/github.com/pangudashu/foolgo/example/www",               //静态文件目录
		ViewPath:   "/home/qinpeng/mygo/src/github.com/pangudashu/foolgo/example/application/views", //模板目录
		Addr:       ":8090",                                                                         //监听地址:端口
		AccessLog:  "/home/qinpeng/log/foolgo/access.log",                                           //请求日志
		ErrorLog:   "/home/qinpeng/log/foolgo/error.log",                                            //错误日志
		RunLog:     "/home/qinpeng/log/foolgo/run.log",                                              //运行日志
		IsGzip:     true,                                                                            //是否开启gzip
		ZipMinSize: 100,                                                                             //gzip压缩起始大小
		Pid:        "/tmp/example.pid",                                                              //进程号保存地址
	}

	server, err := foolgo.NewServer(server_config)
	if err != nil {
		fmt.Println(err)
		return
	}

	//注册控制器
	server.App.RegController(controller_map)

	//添加静态资源压缩类型
	//默认值已有.css .js .html
	static_compress_ext := []string{".txt", ".htm"}
	server.App.AddCompressType(static_compress_ext)

	//Run
	server.Run()
}

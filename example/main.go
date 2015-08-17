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
	server_config := &foolgo.HttpServerConfig{
		Root:       "/home/qinpeng/mygo/src/github.com/pangudashu/foolgo/example/www",               //静态文件目录
		ViewPath:   "/home/qinpeng/mygo/src/github.com/pangudashu/foolgo/example/application/views", //模板目录
		Addr:       ":8090",                                                                         //监听地址:端口
		IsGzip:     true,                                                                            //是否开启gzip
		ZipMinSize: 500,                                                                             //gzip压缩起始大小
		Pid:        "/tmp/example.pid",                                                              //进程号保存地址
	}

	server, err := foolgo.NewServer(server_config)
	if err != nil {
		fmt.Println(err)
		return
	}

	//注册控制器
	server.App.RegController(controller_map)
	//Run
	server.Run()
}

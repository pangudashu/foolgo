package main

import (
	"fmt"
	"github.com/pangudashu/FoolGo"
)

func main() {
	server_config := &foolgo.HttpServerConfig{
		Root:       "/home/qinpeng/mygo/src/github.com/pangudashu/FoolGo/example/www",
		Addr:       ":8090",
		IsGzip:     true,
		ZipMinSize: 500,
		Pid:        "/tmp/example.pid",
	}

	server, err := foolgo.NewServer(server_config)
	if err != nil {
		fmt.Println(err)
		return
	}
	server.Run()
}

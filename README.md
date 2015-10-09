# foolgo
foolgo是一个golang实现的web开发框架,支持热升级、https、自定义路由、rewrite、gzip/deflate压缩等

### 简介
golang相关的web框架已经有很多了，像beego、martini、revel，为什么还要重复造轮子呢？

首先，笔者写这个项目的主要目的很单纯：学习！foolgo参考了很多beego的实现，beego是一个非常不错的项目，提供了很多非常好的范例[@astaxie](https://github.com/astaxie)

其次，golang对web开发已经提供了很丰富的扩展库，通过几行代码即可实现一个server，框架的意义在哪？笔者认为web框架最核心的部分就是路由，即由http->handler的过程。框架最大的意义在于将一些使用频率比较的高的操作进行封装以方便后续的使用。

### 特点
* 支持热升级:不中断服务重启server(抄袭beego ^_^)
* 支持自定义路由
* 支持https
* 支持gzip/deflate压缩
* 支持静态文件

### 使用
	go get github.com/pangudashu/foolgo
	具体使用参考示例

### 示例
Demo目录 ：$GOPATH/src/github.com/pangudashu/foolgo/example
	
	# cd $GOPATH/src/github.com/pangudashu/foolgo/example
	# vim main.go

	//修改server config
	server_config := &foolgo.HttpServerConfig{
	Root:        "{YOUR_GOPATH_DIR}/src/github.com/pangudashu/foolgo/example/www",               //静态文件目录
	ViewPath:    "{YOUR_GOPATH_DIR}/src/github.com/pangudashu/foolgo/example/application/views", //模板目录
	Addr:        ":8090",                                                                        //监听地址:端口
	AccessLog:   "/var/log/foolgo/access.log",                                                   //请求日志
	ErrorLog:    "/var/log/foolgo/error.log",                                                    //错误日志
	RunLog:      "/var/log/foolgo/run.log",                                                      //运行日志
	Compress:    foolgo.COMPRESS_GZIP,                                                           //encoding类型
	CompressMin: 100,                                                                            //encoding最小值
	Pid:         "/tmp/example.pid",                                                             //进程号保存地址
	}

	//编译
	# go install
	//启动
	# $GOPATH/bin/example

	浏览器访问：http://IP:8090/?m=demo.index&id=1234 或 http://IP:8090/demo/1234 (controller中自定义路由)




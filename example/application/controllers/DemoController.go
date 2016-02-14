package controllers

import (
	"github.com/pangudashu/foolgo"
	"time"
)

type DemoController struct {
	foolgo.Controller
}

func (demo *DemoController) IndexAction() {
	demo.Assign("time", time.Now())
	demo.Assign("title", "welcome to foolgo~")
	demo.Assign("id", demo.Param("id"))

	demo.Display()
}

func (demo *DemoController) RegRouter() map[string]interface{} {
	request_map := map[string]interface{}{
		"/demo/:id": "IndexAction",
	}
	return request_map
}

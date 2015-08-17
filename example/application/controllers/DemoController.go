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
	demo.Display()
}

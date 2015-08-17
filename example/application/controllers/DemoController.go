package controllers

import (
	"fmt"
	"github.com/pangudashu/FoolGo"
)

type DemoController struct {
	FoolGo.Controller
}

func (demo *DemoController) IndexAction() {
	fmt.Println("ddd")
}

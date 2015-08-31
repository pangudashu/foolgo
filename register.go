package foolgo

import (
	"fmt"
	"reflect"
)

type Register struct {
	controller_map     map[string]interface{}       //controller_name => controller_type
	router_rewrite_map map[string][]*router_rewrite //url => controller_name
	router_rewrite_key map[string]string
}

var (
	register_instance               *Register
	controller_request_map_register = "RegisterMap"
)

func NewRegister() *Register {
	if register_instance != nil {
		return register_instance
	}
	register_instance = &Register{
		controller_map:     make(map[string]interface{}),
		router_rewrite_map: make(map[string][]*router_rewrite),
		router_rewrite_key: make(map[string]string),
	}

	return register_instance
}

func GetRegister() *Register {
	return register_instance
}

func (this *Register) SetController(controller_name string, controller FGController) error {
	if _, ok := this.controller_map[controller_name]; ok == true {
		logger.RunLog("[Error] conflicting controller name:" + controller_name)
		return fmt.Errorf("%q is existed!", controller_name)
	}

	controller_value := reflect.Indirect(reflect.ValueOf(controller))
	this.controller_map[controller_name] = controller_value.Type()

	return RegRouter(controller, controller_name)
}

func (this *Register) GetRouterWriteMap() map[string][]*router_rewrite {
	return this.router_rewrite_map
}

func (this *Register) GetController(controller_name string) reflect.Type {
	if c, ok := this.controller_map[controller_name]; ok == false {
		return nil
	} else {
		return c.(reflect.Type)
	}
}

func (this *Register) GetAllController() map[string]interface{} {
	return this.controller_map
}

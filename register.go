package foolgo

import (
	//"errors"
	"fmt"
	"reflect"
	//"strings"
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

/*{{{ func NewRegister() *Register
 */
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

/*}}}*/

func GetRegister() *Register {
	return register_instance
}

/*{{{ func (this *Register) SetController(controller_name string,interface{}) error
 */
func (this *Register) SetController(controller_name string, controller FGController) error {
	if _, ok := this.controller_map[controller_name]; ok == true {
		return fmt.Errorf("%q is existed!", controller_name)
	}

	controller_value := reflect.Indirect(reflect.ValueOf(controller))
	this.controller_map[controller_name] = controller_value.Type()

	return RegRouter(controller, controller_name)
	//注册map url
	/*
		url_map := controller.RegRouter()
		if url_map == nil {
			return nil
		}

		for pattern, tb := range url_map {
			switch tb.(type) {
			case string:
				if old, ok := this.router_rewrite_key["GET "+pattern]; ok == true {
					return errors.New("Can't register router:\"GET " + pattern + "\",it had been registed in controller:" + old)
				}
				r := createRewriteRouter(pattern, "GET", controller_name, tb.(string))
				if r == nil {
					return errors.New("Can't register router:" + pattern)
				}
				this.router_rewrite_map["GET"] = append(this.router_rewrite_map["GET"], r)
				this.router_rewrite_key["GET "+pattern] = controller_name
			case map[string]string:
				for method, action := range tb.(map[string]string) {
					method = strings.ToUpper(method)
					if old, ok := this.router_rewrite_key[method+" "+pattern]; ok == true {
						return errors.New("Can't register router:\"" + method + " " + pattern + "\",it had been registed in controller:" + old)
					}
					r := createRewriteRouter(pattern, method, controller_name, action)
					if r == nil {
						return errors.New("Can't register router:\"" + method + " " + pattern + "\"")
					}
					this.router_rewrite_map[method] = append(this.router_rewrite_map[method], r)
					this.router_rewrite_key[method+" "+pattern] = controller_name
				}
			}
		}
	*/
}

/*}}}*/

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

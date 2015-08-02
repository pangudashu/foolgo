package foolgo

import (
	"errors"
	"reflect"
	"strings"
)

type Router struct {
}

type router_rewrite struct {
	urlParts   map[int]map[string]string
	staticNum  int
	method     string
	controller string
	action     string
}

var (
	router_instance    *Router
	router_rewrite_map = make(map[string][]*router_rewrite)
	router_rewrite_key = make(map[string]string)
)

func NewRouter() *Router {
	if router_instance != nil {
		return router_instance
	}
	router_instance = &Router{}

	return router_instance
}

func GetRouter() *Router {
	return router_instance
}

func RegRouter(controller FGController, controller_name string) error {
	url_map := controller.RegRouter()
	if url_map == nil {
		return nil
	}

	for pattern, tb := range url_map {
		switch tb.(type) {
		case string:
			if old, ok := router_rewrite_key["GET "+pattern]; ok == true {
				return errors.New("Can't register router:\"GET " + pattern + "\",it had been registed in controller:" + old)
			}
			r := createRewriteRouter(pattern, "GET", controller_name, tb.(string))
			if r == nil {
				return errors.New("Can't register router:" + pattern)
			}
			router_rewrite_map["GET"] = append(router_rewrite_map["GET"], r)
			router_rewrite_key["GET "+pattern] = controller_name
		case map[string]string:
			for method, action := range tb.(map[string]string) {
				method = strings.ToUpper(method)
				if old, ok := router_rewrite_key[method+" "+pattern]; ok == true {
					return errors.New("Can't register router:\"" + method + " " + pattern + "\",it had been registed in controller:" + old)
				}
				r := createRewriteRouter(pattern, method, controller_name, action)
				if r == nil {
					return errors.New("Can't register router:\"" + method + " " + pattern + "\"")
				}
				router_rewrite_map[method] = append(router_rewrite_map[method], r)
				router_rewrite_key[method+" "+pattern] = controller_name
			}
		}
	}
	return nil
}

func createRewriteRouter(pattern, method, controller, action string) *router_rewrite {
	if pattern == "" || action == "" {
		return nil
	}
	url_parts := strings.Split(strings.Trim(pattern, "/"), "/")

	router := &router_rewrite{
		urlParts:   make(map[int]map[string]string),
		staticNum:  0,
		method:     method,
		controller: controller,
		action:     action,
	}

	for pos, part := range url_parts {
		if part[0:1] == ":" {
			router.urlParts[pos] = map[string]string{
				"name": part[1:],
				"type": "var",
			}
		} else if part[0:1] == "*" {
			router.urlParts[pos] = map[string]string{
				"name": part,
				"type": "",
			}

			break
		} else {
			router.urlParts[pos] = map[string]string{
				"name": part,
				"type": "",
			}
			router.staticNum++
		}
	}
	return router
}

func (this *Router) MatchRewrite(url, method string) (string, string, map[string]string, error) {
	if _, ok := router_rewrite_map[method]; ok == false {
		return "", "", nil, errors.New("No match")
	}

	paths := strings.Split(strings.Trim(url, "/"), "/")
	for _, router := range router_rewrite_map[method] {
		if match_param, matched := this.matchRouter(router, paths); matched == false {
			continue
		} else {
			return router.controller, router.action, match_param, nil
		}
	}

	return "", "", nil, errors.New("No match")
}

func (this *Router) matchRouter(router *router_rewrite, paths []string) (map[string]string, bool) {
	var match_param map[string]string
	var cnt int
	for pos, part := range paths {
		if _, ok := router.urlParts[pos]; ok == false {
			return nil, false
		}
		if router.urlParts[pos]["type"] == "" {
			if router.urlParts[pos]["name"] == "*" {
				if pos < len(paths) {
					if match_param == nil {
						match_param = make(map[string]string)
					}

					param := paths[pos:]
					param_num := len(param) / 2
					for i := 0; i < param_num; i++ {
						match_param[param[i*2]] = param[i*2+1]
					}
				}
				break
			}

			if router.urlParts[pos]["name"] != part {
				return nil, false
			}
			cnt++
		} else if router.urlParts[pos]["type"] == "var" {
			if match_param == nil {
				match_param = make(map[string]string)
			}
			match_param[router.urlParts[pos]["name"]] = part
		}
	}

	if cnt != router.staticNum {
		return nil, false
	}
	return match_param, true
}

/*{{{ func (this *Router) ParseMethod() error
 */
func (this *Router) ParseMethod(method string) (controller_name string, action_name string) {
	method_map := strings.SplitN(method, ".", 2)
	switch len(method_map) {
	case 1:
		controller_name = method_map[0]
	case 2:
		controller_name = method_map[0]
		action_name = method_map[1]
	}
	return controller_name, action_name
}

/*}}}*/

/*{{{ func (this *Router) NewController(controller_name string) (reflect.Value, HTTP_STATUS)
 */
func (this *Router) NewController(controller_name string) (reflect.Value, error) {
	register := GetRegister()

	//m_arr := make([]reflect.Value, 0)
	var controller_instance reflect.Value

	if register == nil {
		//http 500
		return controller_instance, errors.New("Server Error : Can't find \"Register\"")
	}

	controller_type := register.GetController(controller_name)
	if controller_type == nil {
		//http 404
		return controller_instance, errors.New("Warn : Can't find " + controller_name)
	}

	controller_instance = reflect.New(controller_type)
	if false == controller_instance.IsValid() {
		//http 404
		return controller_instance, errors.New("Warn : Can't find " + controller_name)
	}

	return controller_instance, nil
}

/*}}}*/

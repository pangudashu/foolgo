package config

import (
	"fmt"
)

type Config interface {
	Parse(config_path string) (ConfigContainer, error)
}

type ConfigContainer interface {
	String(key string) (string, error)
	Int(key string) (int, error)
	Int64(key string) (int64, error)
	Bool(key string) (bool, error)
	Float(key string) (float64, error)
}

var adapters = make(map[string]Config)

/*{{{ func Register(conf_type string, adapter Config)
 * 注册适配器
 */
func Register(conf_type string, adapter Config) {
	if adapter == nil {
		panic("config: Register adapter is nil")
	}
	if _, ok := adapters[conf_type]; ok {
		panic("config: Register called twice for adapter " + conf_type)
	}
	adapters[conf_type] = adapter
}

/*}}}*/

/*{{{ func GetConfig(conf_type string, config_path string) (ConfigContainer, error)
 */
func GetConfig(conf_type string, config_path string) (ConfigContainer, error) {
	adapter, ok := adapters[conf_type]
	if !ok {
		return nil, fmt.Errorf("config: unknown adaptername %q (forgotten import?)", conf_type)
	}
	return adapter.Parse(config_path)
}

/*}}}*/

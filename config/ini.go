package config

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	empty        = []byte{}
	AnnoPrefix   = []byte{';'}
	ModLeft      = []byte{'['}
	ModRight     = []byte{']'}
	Equal        = []byte{'='}
	SectionSplit = []byte{'.'}
)

type config_node map[string]interface{}

type IniConfig struct {
}

type IniConfigContainer struct {
	data config_node
}

func init() {
	Register("ini", &IniConfig{})
}

/*{{{ func (this *IniConfig) Parse(config_path string) (ConfigContainer, error)
 * 解析配置文件
 */
func (this *IniConfig) Parse(config_path string) (ConfigContainer, error) {
	fp, err := os.Open(config_path)
	defer fp.Close()
	if err != nil {
		return nil, err
	}

	cfg := &IniConfigContainer{
		//make(map[string]map[string]string),
		make(config_node),
	}

	buff := bufio.NewReader(fp)

	var mod string
	for {
		line, _, err := buff.ReadLine()
		if io.EOF == err {
			break
		}
		if bytes.Equal(line, empty) {
			continue
		}
		line = bytes.TrimSpace(line)
		if bytes.HasPrefix(line, AnnoPrefix) {
			continue
		}
		if bytes.HasPrefix(line, ModLeft) && bytes.HasSuffix(line, ModRight) {
			mod = string(line[1 : len(line)-1])
			continue
		}

		if mod == "" {
			continue
		}
		if _, ok := cfg.data[mod]; ok == false {
			cfg.data[mod] = make(config_node)
		}

		vals := bytes.SplitN(line, Equal, 2)
		if len(vals) < 2 {
			continue
		}
		section_val := bytes.SplitN(vals[0], SectionSplit, 2)
		if len(section_val) < 2 {
			cfg.data[mod].(config_node)[string(bytes.TrimSpace(vals[0]))] = string(bytes.TrimSpace(vals[1]))
		} else {
			section := string(bytes.TrimSpace(section_val[0]))
			if _, ok := cfg.data[mod].(config_node)[section]; ok == false {
				cfg.data[mod].(config_node)[section] = make(config_node)
			}

			cfg.data[mod].(config_node)[section].(config_node)[string(bytes.TrimSpace(section_val[1]))] = string(bytes.TrimSpace(vals[1]))
		}
	}

	return cfg, nil
}

/*}}}*/

/*{{{ func (container *IniConfigContainer) getData(key string) (string, error)
 */
func (container *IniConfigContainer) getData(key string) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	key_arr := strings.Split(key, ":")
	if len(key_arr) < 2 {
		return "", fmt.Errorf("config: error usage!you should get a config value like this : \"Module:Section.Variable\"")
	}
	mod := key_arr[0]
	if _, ok := container.data[mod]; ok == false {
		return "", fmt.Errorf("config: there is no Config-Module %q", mod)
	}

	section_arr := strings.Split(key_arr[1], ".")
	if len(section_arr) == 1 {
		if d, ok := container.data[mod].(config_node)[section_arr[0]]; ok == false {
			return "", fmt.Errorf("config: there is no config-variable %q", section_arr[0])
		} else {
			return d.(string), nil
		}
	} else if len(section_arr) == 2 {
		section := section_arr[0]
		variable := section_arr[1]
		if _, ok := container.data[mod].(config_node)[section]; ok == false {
			return "", fmt.Errorf("config: there is no config-section %q", section)
		}
		if d, ok := container.data[mod].(config_node)[section].(config_node)[variable]; ok == false {
			return "", fmt.Errorf("config: there is no config-variable %q", variable)
		} else {
			return d.(string), nil
		}
	} else {
		return "", fmt.Errorf("config: error usage!you should get a config value like this : \"Module:Section.Variable\"")
	}
}

/*}}}*/

/*{{{* func (container *IniConfigContainer) String(key string) (string, error)
 */
func (container *IniConfigContainer) String(key string) (string, error) {
	return container.getData(key)
}

/*}}}*/

/*{{{* func (container *IniConfigContainer) Int(key string) (int, error)
 */
func (container *IniConfigContainer) Int(key string) (int, error) {
	val, err := container.getData(key)

	if err != nil {
		return 0, err
	}
	return strconv.Atoi(val)
}

/*}}}*/

/*{{{ func (container *IniConfigContainer) Int64(key string) (int64, error)
 */
func (container *IniConfigContainer) Int64(key string) (int64, error) {
	val, err := container.getData(key)

	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(val, 10, 64)
}

/*}}}*/

/*{{{ func (container *IniConfigContainer) Bool(key string) (bool, error)
 */
func (container *IniConfigContainer) Bool(key string) (bool, error) {
	val, err := container.getData(key)

	if err != nil {
		return false, err
	}

	return strconv.ParseBool(val)
}

/*}}}*/

/*{{{ func (container *IniConfigContainer) Float(key string) (float64, error)
 */
func (container *IniConfigContainer) Float(key string) (float64, error) {
	val, err := container.getData(key)

	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(val, 64)
}

/*}}}*/

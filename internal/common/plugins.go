package common

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrInvalidePluginOption = errors.New("invalid plugin option")
	ErrInvalidOptionValue   = errors.New("invalid option value")
)

type Plugin interface {
	Name() string
	Run(Context, PluginOptions, chan interface{}) error
}

func Plugins(t reflect.Type) []interface{} {
	plugins := [](interface{}){}

	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if !strings.HasPrefix(m.Name, "Register") {
			continue
		}

		x := reflect.New(t).Elem().Interface()
		y := m.Func.Call([]reflect.Value{reflect.ValueOf(x)})[0].Interface()
		//fmt.Println("plugin found:", reflect.TypeOf(y), y)

		plugins = append(plugins, y)
	}
	return plugins
}

// Runner is a function type for plugable checkers and exporters.
type Runner func(Context, PluginOptions, chan interface{}) error

// PluginOptions is used for CheckerOptions and ExporterOptions.
type PluginOptions map[string][]string

// BuildPluginOptions converts given string s into a map of PluginOptions.
// The string contains all options of all plugins at once. Option entries
// are separated by ';' and Each option entry consists of three parts with
// delimiter ':'. The first part is the name of plugin, the second is the
// name of the option, and the third is the value of the option. The value
// part is a comma or space separated list of values even if there is only
// one value (the plugins may understand them):
/*
	"plugin:opt1_key:opt1_value1,opt1_value2;plugin:opt2_key:opt2_value"
*/
func BuildPluginOptions(s string) (map[string]PluginOptions, error) {
	options := map[string]PluginOptions{}
	if len(s) > 1 {
		for _, e := range strings.Split(s, ";") {
			t := strings.SplitN(e, ":", 3)
			if len(t) == 3 {
				_, found := options[t[0]]
				if found {
					options[t[0]][t[1]] = StringValues(t[2])
				} else {
					options[t[0]] = PluginOptions{t[1]: StringValues(t[2])}
				}
			} else {
				return nil, fmt.Errorf("%w: %v", ErrInvalidePluginOption, e)
			}
		}
	}
	return options, nil
}

// StringValues creates and returns a list of strings from the given string s.
// If the string contains any comma, commas will be used as delimiter and
// all leading/trailing spaces of each substring will be removed. When there
// is no comma on the string, (merged) white-spaces are used as a delimiter.
/*
    "hey, bulldog " --> ["hey", "dog"]
	" oh  darling " --> ["oh", "darling"]
*/
func StringValues(s string) []string {
	ret := []string{}
	if csv := strings.Split(s, ","); len(csv) > 1 {
		for _, v := range csv {
			ret = append(ret, strings.TrimSpace(v))
		}
	} else {
		ret = strings.Fields(s)
	}
	return ret
}

func (o *PluginOptions) GetValuesOr(key string, def []string) []string {
	ret := (*o)[key]
	if len(ret) > 0 {
		return ret
	}
	return def
}

func (o *PluginOptions) GetValueOr(key, def string) string {
	values := o.GetValuesOr(key, []string{})
	if len(values) > 0 {
		return values[0]
	}
	return def
}

func (o *PluginOptions) GetIntegerOr(key string, def int) (int, error) {
	value := o.GetValueOr(key, strconv.Itoa(def))
	if ret, err := strconv.Atoi(value); err != nil {
		return def, err
	} else {
		return ret, nil
	}
}

// utilities ---

// Contains checks if the given list has the given string item.
func Contains(list []string, item string) bool {
	for _, e := range list {
		if e == item {
			return true
		}
	}
	return false
}

package common

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	// ErrInvalidePluginOption indicates that the given option name is unknown
	// from the perspective of the plugin. It could also indicate a typos.
	ErrInvalidePluginOption = errors.New("invalid plugin option")
	// ErrInvalidOptionValue indicates that the given option value is invalid.
	// It could used for the situation that numeric conversion is not possible.
	ErrInvalidOptionValue = errors.New("invalid option value")
)

// Plugin is an interface for all bogo plugins.
type Plugin interface {
	// Name returns the name of the plugin.
	Name() string
	// Run starts the plugin's runner in a goroutine and returns the status of
	// execution. The runner will run forever until the given context is
	// canceled.
	Run(Context, PluginOptions, chan interface{}) error
}

// Plugins returns the plugin list for the given type. All register function
// of the type which is prefixed with "Register" will executed during the
// registration.
func Plugins(t reflect.Type) []interface{} {
	plugins := [](interface{}){}

	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if !strings.HasPrefix(m.Name, "Register") {
			continue
		}

		x := reflect.New(t).Elem().Interface()
		y := m.Func.Call([]reflect.Value{reflect.ValueOf(x)})[0].Interface()
		// fmt.Println("plugin found:", reflect.TypeOf(y), y)

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

	if s == "" {
		return options, nil
	}

	for _, e := range strings.Split(s, ";") {
		t := strings.SplitN(e, ":", 3) // nolint
		if len(t) != 3 {               // nolint
			return nil, fmt.Errorf("%w: %v", ErrInvalidePluginOption, e)
		}

		_, pluginExists := options[t[0]]
		if pluginExists {
			options[t[0]][t[1]] = StringValues(t[2])
		} else {
			options[t[0]] = PluginOptions{t[1]: StringValues(t[2])}
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
	if csv := strings.Split(s, ","); len(csv) > 1 {
		ret := []string{}
		for _, v := range csv {
			ret = append(ret, strings.TrimSpace(v))
		}

		return ret
	}

	return strings.Fields(s)
}

// GetValuesOr returns the option values for the given key from the options.
// If no option values found, it will returns the def value.
func (o *PluginOptions) GetValuesOr(key string, def []string) []string {
	if ret := (*o)[key]; len(ret) > 0 {
		return ret
	}

	return def
}

// GetValueOr returns the option value for the given key from the options as
// string. If no option value found, it will returns the def string.
//
// Note that it just returns the first value if the option values are more than
// one.
func (o *PluginOptions) GetValueOr(key, def string) string {
	if values := o.GetValuesOr(key, []string{}); len(values) > 0 {
		return values[0]
	}

	return def
}

// GetIntegerOr returns the option value for the given key from the options as
// int. If no option value found, it will returns the def value with nil error.
//
// Note that if the found value is not properly converted to integer, it will
// returns the default value with error.
func (o *PluginOptions) GetIntegerOr(key string, def int) (int, error) {
	value := o.GetValueOr(key, strconv.Itoa(def))
	if ret, err := strconv.Atoi(value); err != nil {
		return def, err
	} else { // nolint
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

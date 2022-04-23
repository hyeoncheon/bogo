package common

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPlugins(t *testing.T) {
	r := require.New(t)

	plugins := Plugins(reflect.TypeOf(&Plugger{}))
	r.IsType([]interface{}{}, plugins)
	r.NotEmpty(plugins)
	r.Equal(2, len(plugins))

	e1 := plugins[0].(Plugin)
	r.IsType(&Plugger{}, e1)
	r.Equal("dummy", e1.Name())
	r.Nil(e1.Run(nil, nil, nil))

	e2 := plugins[1].(Plugin)
	r.IsType(&Plugger{}, e2)
	r.Equal("mummy", e2.Name())
	r.Error(e2.Run(nil, nil, nil))
}

func TestBuildPluginOptions(t *testing.T) {
	r := require.New(t)

	tests := []struct {
		in  string
		out map[string]PluginOptions
		e   error
	}{
		{
			in: "heartbeat:interval:1;ping:targets:www.google.com,ns.kornet.net;ping:check_interval:20",
			out: map[string]PluginOptions{
				"heartbeat": {"interval": []string{"1"}},
				"ping": {
					"targets":        []string{"www.google.com", "ns.kornet.net"},
					"check_interval": []string{"20"},
				},
			},
			e: nil,
		},
		{
			in: "heartbeat:interval:1",
			out: map[string]PluginOptions{
				"heartbeat": {"interval": []string{"1"}},
			},
			e: nil,
		},
		{
			in: "heartbeat:interval:1920:1080",
			out: map[string]PluginOptions{
				"heartbeat": {"interval": []string{"1920:1080"}},
			},
			e: nil,
		},
		{
			in:  "heartbeat:interval",
			out: nil,
			e:   fmt.Errorf("%w: heartbeat:interval", ErrInvalidePluginOption),
		},
		{
			in:  "heartbeat",
			out: nil,
			e:   fmt.Errorf("%w: heartbeat", ErrInvalidePluginOption),
		},
		{
			in:  "",
			out: map[string]PluginOptions{},
			e:   nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			out, err := BuildPluginOptions(tc.in)
			r.Equal(tc.e, err)
			r.EqualValues(tc.out, out)
		})
	}
}

func TestStringValues(t *testing.T) {
	r := require.New(t)

	tests := []struct {
		name string
		in   string
		out  []string
	}{
		{
			name: "comma",
			in:   "hey,bulldog",
			out:  []string{"hey", "bulldog"},
		},
		{
			name: "comma_space",
			in:   "hey, bulldog ",
			out:  []string{"hey", "bulldog"},
		},
		{
			name: "space_comma_space",
			in:   " hey , bulldog",
			out:  []string{"hey", "bulldog"},
		},
		{
			name: "comma_spaces",
			in:   "  hey,  bulldog  ",
			out:  []string{"hey", "bulldog"},
		},
		{
			name: "commas_space",
			in:   "hey,, bulldog",
			out:  []string{"hey", "", "bulldog"},
		},
		{
			name: "space",
			in:   "hey bulldog",
			out:  []string{"hey", "bulldog"},
		},
		{
			name: "spaces",
			in:   "  hey  bulldog  ",
			out:  []string{"hey", "bulldog"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := StringValues(tc.in)
			r.EqualValues(tc.out, out)
		})
	}
}

func TestContains(t *testing.T) {
	r := require.New(t)

	list := []string{"hey", "bulldog"}
	r.True(Contains(list, "hey"))
	r.True(Contains(list, "bulldog"))
	r.False(Contains(list, "hotdog"))
}

func TestGetValue(t *testing.T) {
	r := require.New(t)

	opts := PluginOptions{
		"value1":  []string{"one"},
		"value2":  []string{"first", "second"},
		"value3":  []string{"6090"},
		"invalid": []string{"bogo"},
	}

	val1 := opts.GetValueOr("value1", "default")
	r.Equal("one", val1)

	val1 = opts.GetValueOr("none", "default")
	r.Equal("default", val1)

	val2 := opts.GetValuesOr("value2", []string{"default"})
	r.Equal([]string{"first", "second"}, val2)

	val2 = opts.GetValuesOr("none", []string{"default"})
	r.Equal([]string{"default"}, val2)

	val3, err := opts.GetIntegerOr("value3", 8080)
	r.NoError(err)
	r.Equal(6090, val3)

	val3, err = opts.GetIntegerOr("none", 8080)
	r.NoError(err)
	r.Equal(8080, val3)

	val3, err = opts.GetIntegerOr("invalid", 8080)
	r.Error(err)
	r.Equal(8080, val3)
}

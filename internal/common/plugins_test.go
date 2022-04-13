package common

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

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
			e:   errors.New("invalid plugin option 'heartbeat:interval'"),
		},
		{
			in:  "heartbeat",
			out: nil,
			e:   errors.New("invalid plugin option 'heartbeat'"),
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

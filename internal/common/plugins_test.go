package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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

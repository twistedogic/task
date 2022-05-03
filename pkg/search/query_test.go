package search

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_ParseQuery(t *testing.T) {
	cases := map[string]struct {
		input string
		want  Query
	}{
		"base": {
			input: "term search f:.go r:some_git",
			want: Query{
				Repo: "some_git",
				File: ".go",
				Term: "term search",
			},
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			got, err := ParseQuery(strings.Split(tc.input, " ")...)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

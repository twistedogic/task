package cheat

import (
	"testing"
)

func TestGetAnswer(t *testing.T) {
	cases := []struct {
		lang     string
		query    []string
		quiet    bool
		page     int
		expected string
	}{
		{
			"go", []string{"how", "to", "save", "file"}, true, 0, "",
		},
	}
	for _, test := range cases {
		out, err := GetAnswer(test.lang, test.page, test.quiet, test.query...)
		if err != nil {
			t.Error(err)
		}
		if out == "" {
			t.Fail()
		}
	}
}

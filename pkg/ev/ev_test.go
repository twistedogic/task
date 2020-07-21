package ev

import (
	"testing"
)

func Test_GetBreakeven(t *testing.T) {
	cases := map[string]struct {
		input    []string
		want     float64
		hasError bool
	}{
		"base":  {[]string{"100", "0.5"}, 100.0, false},
		"empty": {[]string{}, 0.0, true},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			got, err := GetBreakeven(tc.input)
			if (err != nil) != tc.hasError {
				t.Fatal(err)
			}
			if !tc.hasError {
				if got != tc.want {
					t.Fatalf("want: %f, got: %f", tc.want, got)
				}
			}
		})
	}
}

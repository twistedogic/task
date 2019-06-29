package objective

import (
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/twistedogic/task/pkg/okr/keyresult"
)

func TestObjectiveLoad(t *testing.T) {
	cases := map[string]struct {
		input    []byte
		want     *Objective
		hasError bool
	}{
		"base": {
			[]byte(`{
				"name":"test",
				"key_results": [],
				"board": "board",
				"last_update": "2001-01-02T00:00:00Z",
				"target": "2001-01-03T00:00:00Z"
			}`),
			&Objective{
				Name:       "test",
				KeyResults: []*keyresult.KeyResult{},
				Index:      map[string]*keyresult.KeyResult{},
				Board:      "board",
				LastUpdate: time.Date(2001, time.January, 2, 0, 0, 0, 0, time.UTC),
				Target:     time.Date(2001, time.January, 3, 0, 0, 0, 0, time.UTC),
				Done:       false,
				Score:      0.0,
			},
			false,
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			got, err := Load(tc.input)
			hasError := err != nil
			if hasError != tc.hasError {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestByTargetDate(t *testing.T) {
	objs := []*Objective{
		&Objective{
			Name:       "test",
			KeyResults: []*keyresult.KeyResult{},
			Index:      map[string]*keyresult.KeyResult{},
			Board:      "board",
			LastUpdate: time.Date(2001, time.January, 2, 0, 0, 0, 0, time.UTC),
			Target:     time.Date(2001, time.January, 3, 0, 0, 0, 0, time.UTC),
			Done:       false,
			Score:      0.0,
		},
		&Objective{
			Name:       "test",
			KeyResults: []*keyresult.KeyResult{},
			Index:      map[string]*keyresult.KeyResult{},
			Board:      "board",
			LastUpdate: time.Date(2001, time.January, 2, 0, 0, 0, 0, time.UTC),
			Target:     time.Date(2000, time.January, 3, 0, 0, 0, 0, time.UTC),
			Done:       false,
			Score:      0.0,
		},
	}
	sort.Sort(ByTargetDate(objs))
	for i, obj := range objs {
		if i == 0 {
			continue
		}
		if obj.Target.Before(objs[i-1].Target) {
			t.Fail()
		}
	}
}

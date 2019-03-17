package item

import (
	"reflect"
	"testing"
	"time"

	"github.com/bouk/monkey"
)

func TestNewTask(t *testing.T) {
	mockNow := time.Date(1974, time.May, 19, 1, 2, 3, 4, time.UTC)
	patch := monkey.Patch(time.Now, func() time.Time { return mockNow })
	defer patch.Unpatch()
	expect := &Item{
		ID:          uint(0),
		Description: "ok",
		Board:       "board1",
		CreateDate:  mockNow,
		LastUpdate:  mockNow,
	}
	got := NewItem(uint(0), TASK, "ok", "board1")
	if !reflect.DeepEqual(expect, got) {
		t.Errorf("Expect: %+v, Got: %+v\n", expect, got)
	}
}

func TestSortByUpdateDateDesc(t *testing.T) {
	format := "2006-01-02"
	dates := []string{"2018-01-02", "2017-12-31", "2018-01-03"}
	expect := []string{"2018-01-03", "2018-01-02", "2017-12-31"}
	input := make([]*Item, len(dates))
	for i, v := range dates {
		t, _ := time.Parse(format, v)
		input[i] = &Item{
			LastUpdate: t,
		}
	}
	SortBy(byUpdateDateDesc).Sort(input)
	for i, v := range input {
		got := v.LastUpdate.Format(format)
		if expect[i] != got {
			t.Errorf("Expect: %s, Got: %s\n", expect[i], got)
		}
	}
}

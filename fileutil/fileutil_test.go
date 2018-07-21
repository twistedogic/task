package fileutil

import "testing"

func TestHome(t *testing.T) {
	_, err := Home()
	if err != nil {
		t.Fail()
	}
}

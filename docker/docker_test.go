package docker

import "testing"

func TestRunTask(t *testing.T) {
	cases := []struct {
		input  []string
		expect string
	}{
		{[]string{"busybox", "echo", "hi"}, "hi\n"},
	}
	for _, test := range cases {
		out, err := runTask(test.input...)
		if err != nil {
			t.Error(err)
		}
		if string(out) != test.expect {
			t.Errorf("Expect: %s, Got: %s", test.expect, string(out))
		}
	}
}

func TestIsDockerRunning(t *testing.T) {
	t.Log(IsDockerRunning())
}

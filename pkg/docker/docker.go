package docker

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

var (
	DOCKER_NOT_RUNNING = errors.New("docker is not running")
)

func RunTask(args ...string) ([]byte, error) {
	if !IsDockerRunning() {
		return nil, DOCKER_NOT_RUNNING
	}
	root := "docker"
	arg := append([]string{"run", "--rm"}, args...)
	return exec.Command(root, arg...).CombinedOutput()
}

func IsDockerRunning() bool {
	root := "docker"
	cmd := exec.Command(root, "ps")
	return cmd.Run() == nil
}

func StartDocker() error {
	return exec.Command("open", "--background", "-a", "Docker").Run()
}

func isValidLangImage(lang string) error {
	root := "docker"
	b, err := exec.Command(root, "images").CombinedOutput()
	if err != nil {
		return err
	}
	for i, line := range strings.Split(string(b), "\n") {
		if i == 0 || len(line) == 0 {
			continue
		}
		if image := strings.Fields(line)[0]; strings.HasPrefix(image, lang) {
			return nil
		}
	}
	return fmt.Errorf("%s: %s", INVALID_LANG, lang)
}

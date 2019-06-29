package docker

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

const (
	DOCKER_NOT_RUNNING = "docker is not running"
)

func runTask(args ...string) ([]byte, error) {
	if !IsDockerRunning() {
		return []byte{}, fmt.Errorf("%s", DOCKER_NOT_RUNNING)
	}
	root := "docker"
	arg := append([]string{"run", "--rm"}, args...)
	cmd := exec.Command(root, arg...)
	return cmd.CombinedOutput()
}

func IsDockerRunning() bool {
	root := "docker"
	cmd := exec.Command(root, "ps")
	return cmd.Run() == nil
}

func StartDocker() error {
	cmd := exec.Command("open", "--background", "-a", "Docker")
	return cmd.Run()
}

func ListImages() (map[string]bool, error) {
	out := make(map[string]bool)
	root := "docker"
	cmd := exec.Command(root, "images")
	b, err := cmd.CombinedOutput()
	if err != nil {
		return out, err
	}
	lines := strings.Split(string(b), "\n")
	for i, line := range lines {
		if i == 0 || len(line) == 0 {
			continue
		}
		image := strings.Fields(line)[0]
		out[image] = true
	}
	return out, nil
}

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "run container for task",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := runTask(args...)
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(string(out))
		}
	},
}

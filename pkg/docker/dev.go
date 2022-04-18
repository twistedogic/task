package docker

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path"

	"github.com/twistedogic/task/pkg/fileutil"
)

const (
	BASE_PORT    = 3000
	MODE         = 0755
	INVALID_LANG = "invalid language choice"
)

func isPortAvailable(p int) bool {
	port := fmt.Sprintf(":%d", p)
	lis, err := net.Listen("tcp", port)
	if lis != nil {
		lis.Close()
	}
	return err == nil
}

func AssignPort() int {
	var i int
	for {
		port := BASE_PORT + i
		if isPortAvailable(port) {
			return port
		}
		i += 1
	}
	return -1
}

func concat(s ...[]string) []string {
	out := []string{}
	for _, v := range s {
		out = append(out, v...)
	}
	return out
}

func RunDevEnv(lang string) error {
	port := AssignPort()
	u, err := user.Current()
	if err != nil {
		return err
	}
	if !IsDockerRunning() {
		return fmt.Errorf("%s", DOCKER_NOT_RUNNING)
	}
	if err := isValidLangImage(lang); err != nil {
		return err
	}
	home, err := fileutil.Home()
	if err != nil {
		return err
	}
	workspace := path.Join(home, "Dev", fmt.Sprintf("%s_project", lang))
	if err := fileutil.CreateFolder(workspace, MODE); err != nil {
		return err
	}
	run := []string{"run", "-it", "--rm"}
	workspaceVolume := []string{"-v", fmt.Sprintf("%s:/root/workspace", workspace)}
	sshVolume := []string{"-v", fmt.Sprintf("%s/.ssh:/root/.ssh:ro", u.HomeDir)}
	configVolume := []string{"-v", fmt.Sprintf("%s/.config:/root/.config", u.HomeDir)}
	ports := []string{"-p", fmt.Sprintf("%d:%d", port, port)}
	container := []string{fmt.Sprintf("%sbox", lang), "tmux"}
	args := concat(run, workspaceVolume, sshVolume, configVolume, ports, container)
	cmd := exec.Command("docker", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

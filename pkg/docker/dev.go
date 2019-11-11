package docker

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path"

	"github.com/spf13/cobra"
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

func runDevEnv(lang string) error {
	port := AssignPort()
	u, err := user.Current()
	if err != nil {
		return err
	}
	if !IsDockerRunning() {
		return fmt.Errorf("%s", DOCKER_NOT_RUNNING)
	}
	box := fmt.Sprintf("%sbox", lang)
	images, err := ListImages()
	if err != nil {
		return err
	}
	if _, ok := images[box]; !ok {
		return fmt.Errorf("%s: %s", INVALID_LANG, lang)
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
	ports := []string{"-p", fmt.Sprintf("%d:%d", port, port)}
	container := []string{fmt.Sprintf("%sbox", lang), "tmux"}
	args := concat(run, workspaceVolume, sshVolume, ports, container)
	cmd := exec.Command("docker", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

var DevCmd = &cobra.Command{
	Use:   "dev",
	Short: "run container for development",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("not language choice entered")
		}
		lang := args[0]
		if err := runDevEnv(lang); err != nil {
			log.Fatal(err)
		}
	},
}

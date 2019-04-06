package docker

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/cobra"

	"github.com/twistedogic/task/fileutil"
)

const (
	MODE         = 0755
	INVALID_LANG = "invalid language choice"
)

func concat(s ...[]string) []string {
	out := []string{}
	for _, v := range s {
		out = append(out, v...)
	}
	return out
}

func runDevEnv(lang string) error {
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
	volume := []string{"-v", fmt.Sprintf("%s:/root/workspace", workspace)}
	port := []string{"-p", "3000:3000", "-p", "8000:8000"}
	container := []string{fmt.Sprintf("%sbox", lang), "tmux"}
	args := concat(run, volume, port, container)
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

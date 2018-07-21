package note

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/spf13/cobra"

	"github.com/twistedogic/task/fileutil"
)

const (
	FOLDER_MODE = 0700
	FILE_MODE   = 0644
	FORMAT      = "2006-01-02"
)

func StartNote(basePath string, static bool) error {
	home, err := fileutil.Home()
	if err != nil {
		return err
	}
	base := path.Join(home, basePath)
	if err := fileutil.CreateFolder(base, FOLDER_MODE); err != nil {
		return err
	}
	filename := fmt.Sprintf("%s.md", time.Now().Format(FORMAT))
	if static {
		filename = fmt.Sprintf("%s.md", basePath)
	}
	file := path.Join(base, filename)
	if err := fileutil.CreateFile(file, FILE_MODE); err != nil {
		return err
	}
	cmd := exec.Command("vim", file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func createCmd(name, desc string, static bool) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: desc,
		Run: func(cmd *cobra.Command, args []string) {
			if err := StartNote(cmd.Use, static); err != nil {
				log.Fatal(err)
			}
		},
	}
}

var IssueCmd = createCmd("issue", "issue log", true)
var TILCmd = createCmd("til", "Today I Learn", false)
var IdeaCmd = createCmd("idea", "idea log", true)
var FactCmd = createCmd("fact", "fact to check", true)
var BlogCmd = createCmd("blog", "blog", false)

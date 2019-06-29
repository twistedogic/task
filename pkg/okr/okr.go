package okr

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/twistedogic/task/pkg/fileutil"
	"github.com/twistedogic/task/pkg/store/database"
)

const dbName = "okr.db"

func Setup(name string) (*Store, error) {
	dir, err := fileutil.Home()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, ".task", dbName)
	base := filepath.Dir(path)
	if err := fileutil.CreateFolder(base, 0755); err != nil {
		return nil, err
	}
	db, err := database.New(path)
	if err != nil {
		return nil, err
	}
	return New(db)
}

func List(cmd *cobra.Command, args []string) {
	store, err := Setup(dbName)
	if err != nil {
		log.Fatal(err)
	}
	for _, o := range store.List() {
		fmt.Println(o)
	}
}

func Add(cmd *cobra.Command, args []string) {
	db, err := Setup(dbName)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Add(); err != nil {
		log.Fatal(err)
	}
}

func Edit(cmd *cobra.Command, args []string) {
	db, err := Setup(dbName)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Edit(); err != nil {
		log.Fatal(err)
	}
}

func Update(cmd *cobra.Command, args []string) {
	db, err := Setup(dbName)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Update(); err != nil {
		log.Fatal(err)
	}
}

var RunCmd = &cobra.Command{
	Use:   "okr",
	Short: "Objective Key Result",
	Run:   List,
}

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add Objective Key Result",
	Run:   Add,
}

var EditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit Objective Key Result",
	Run:   Edit,
}

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update KRs",
	Run:   Update,
}

func init() {
	RunCmd.AddCommand(AddCmd, EditCmd, UpdateCmd)
}

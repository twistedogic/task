package okr

import (
	"fmt"
	"log"
	"path/filepath"

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

func List() {
	store, err := Setup(dbName)
	if err != nil {
		log.Fatal(err)
	}
	for _, o := range store.List() {
		fmt.Println(o)
	}
}

func Add() {
	db, err := Setup(dbName)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Add(); err != nil {
		log.Fatal(err)
	}
}

func Edit() {
	db, err := Setup(dbName)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Edit(); err != nil {
		log.Fatal(err)
	}
}

func Update() {
	db, err := Setup(dbName)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Update(); err != nil {
		log.Fatal(err)
	}
}

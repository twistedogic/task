package taskbook

import "github.com/twistedogic/task/taskbook/item"

type Storage interface {
	Save() error
	Load() *Taskbook
}

type Index interface {
	Get(uint) *item.Item
	Search(string) []*item.Item
	Update(*item.Item)
}

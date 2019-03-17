package taskbook

import "github.com/twistedogic/task/taskbook/item"

type Taskbook struct {
	Data   []*item.Item
	Index  Index
	LastID uint
	Store  Storage
}

func NewTaskbook(store Storage) *Taskbook {
	t := &Taskbook{Store: store}
	t.Load()
	return t
}

func (t *Taskbook) Load() {
	t = t.Store.Load()
}

func (t *Taskbook) GetItem(id uint) *item.Item {
	return t.Index.Get(id)
}

func (t *Taskbook) AddItem(desc string, itemType item.ItemType, board ...string) {
	ID := t.LastID + uint(1)
	item := item.NewItem(ID, itemType, desc, board...)
	t.Data = append(t.Data, item)
	t.Index.Update(item)
}

func (t *Taskbook) ListItem() {
	groups := item.GroupBy(item.GroupByBoard).Group(t.Data)
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
		item.SortBy(item.ByID).Sort(groups[k])
	}
}

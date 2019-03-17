package item

import "sort"

type SortBy func(t1, t2 *Item) bool

func (s SortBy) Sort(items []*Item) {
	ps := &ItemSorter{
		items: items,
		by:    s,
	}
	sort.Sort(ps)
}

type ItemSorter struct {
	items []*Item
	by    func(t1, t2 *Item) bool
}

func (i *ItemSorter) Len() int {
	return len(i.items)
}

func (i *ItemSorter) Swap(x, y int) {
	i.items[x], i.items[y] = i.items[y], i.items[x]
}

func (i *ItemSorter) Less(x, y int) bool {
	return i.by(i.items[x], i.items[y])
}

var ByUpdateDateDesc = func(t1, t2 *Item) bool {
	if t1.LastUpdate.Before(t2.LastUpdate) {
		return false
	}
	return true
}

var ByID = func(i1, i2 *Item) bool {
	if i1.ID < i2.ID {
		return true
	}
	return false
}

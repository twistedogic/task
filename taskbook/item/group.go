package item

type GroupBy func(i *Item) string

func (g GroupBy) Group(items []*Item) map[string][]*Item {
	groups := make(map[string][]*Item)
	for _, v := range items {
		key := g(v)
		groups[key] = append(groups[key], v)
	}
	return groups
}

var GroupByBoard = func(i *Item) string {
	return i.Board
}

var GroupByType = func(i *Item) string {
	return typeMapping[i.Type]
}

var GroupByPriority = func(i *Item) string {
	return priorityMapping[i.Priority]
}

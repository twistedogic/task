package item

import (
	"time"
)

type ItemType uint
type Priority uint

const (
	TASK ItemType = iota
	NOTE
	IDEA
)

var typeMapping = map[ItemType]string{
	TASK: "task",
	NOTE: "note",
	IDEA: "idea",
}

const (
	HIGH Priority = iota
	MEDIUM
	LOW
)

var priorityMapping = map[Priority]string{
	HIGH:   "high",
	MEDIUM: "medium",
	LOW:    "low",
}

type Item struct {
	ID          uint
	CreateDate  time.Time
	LastUpdate  time.Time
	Description string
	Priority    Priority
	Type        ItemType
	IsComplete  bool
	Board       string
	Assignee    string
}

func NewItem(id uint, itemType ItemType, desc string, board ...string) *Item {
	var b string
	if len(board) > 0 {
		b = board[0]
	}
	return &Item{
		ID:          id,
		CreateDate:  time.Now(),
		LastUpdate:  time.Now(),
		Type:        itemType,
		Description: desc,
		Board:       b,
		IsComplete:  false,
	}
}

func (i *Item) update() {
	i.LastUpdate = time.Now()
}

func (i *Item) UpdateDescription(desc string) {
	i.Description = desc
	i.update()
}

func (i *Item) UpdateBoards(board string) {
	i.Board = board
	i.update()
}

func (i *Item) UpdateAssignee(assignee string) {
	i.Assignee = assignee
	i.update()
}

package search

import (
	"github.com/charmbracelet/bubbles/list"
)

type Match struct {
	Line, Offset, Length int
}

type Result struct {
	Repo, File, Commit, Link string
	Content                  []byte
	Matches                  []Match
}

// implement list.Item
func (r Result) FilterValue() string { return r.Repo + " " + r.File }

// implement list.DefaultItem
func (r Result) Title() string { return r.Repo }

// implement list.DefaultItem
func (r Result) Description() string { return r.File + " " + r.Commit }

type Results []Result

func (r Results) Items() []list.Item {
	items := make([]list.Item, len(r))
	for i, res := range r {
		items[i] = res
	}
	return items
}

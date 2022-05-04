package search

import (
	"fmt"
	"strings"

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

func (r Result) preview() string {
	m := r.Matches[0]
	lines := strings.Split(string(r.Content), "\n")
	return strings.TrimSpace(lines[m.Line])
}

// implement list.Item
func (r Result) FilterValue() string { return r.Repo + " " + r.File }

// implement list.DefaultItem
func (r Result) Title() string { return r.Repo + " " + r.File }

// implement list.DefaultItem
func (r Result) Description() string { return r.preview() }

func (r Result) QueryString() string { return fmt.Sprintf("r:%s f:%s", r.Repo, r.File) }

type Results []Result

func (r Results) Items() []list.Item {
	items := make([]list.Item, len(r))
	for i, res := range r {
		items[i] = res
	}
	return items
}

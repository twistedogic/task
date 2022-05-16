package search

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

type Match struct {
	Line       string
	LineNumber int
}

type Result struct {
	Source                   SourceKey
	Repo, File, Commit, Link string
	Content                  []byte
	Matches                  []Match
}

func (r Result) preview() string {
	if len(r.Matches) == 0 {
		return "no match"
	}
	return strings.TrimSpace(r.Matches[0].Line)
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

package search

import (
	"fmt"
	"strings"
)

type Query struct {
	Repo, Term, File string
}

func (q Query) Validate() error {
	if len(q.Repo+q.Term+q.File) == 0 {
		return fmt.Errorf("query is empty")
	}
	return nil
}

func ParseQuery(args ...string) (Query, error) {
	repo, file, terms := "", "", make([]string, 0, len(args))
	for _, s := range args {
		switch {
		case strings.HasPrefix(s, "f:"):
			file = strings.TrimPrefix(s, "f:")
		case strings.HasPrefix(s, "r:"):
			repo = strings.TrimPrefix(s, "r:")
		default:
			terms = append(terms, s)
		}
	}
	q := Query{File: file, Repo: repo, Term: strings.Join(terms, " ")}
	return q, q.Validate()
}

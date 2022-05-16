package search

import (
	"context"
)

type SourceKey string

type Searcher interface {
	Search(context.Context, Query) ([]Result, error)
}

type Downloader interface {
	Download(context.Context, *Result) error
}

type Source interface {
	Searcher
	Downloader
	IsSource(SourceKey) bool
}

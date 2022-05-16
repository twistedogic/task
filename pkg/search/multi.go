package search

import (
	"context"
	"fmt"
)

type MultiSource struct {
	srcs []Source
}

func NewMultiSource(srcs ...Source) MultiSource {
	return MultiSource{srcs: srcs}
}

func (m MultiSource) Search(ctx context.Context, q Query) ([]Result, error) {
	var err error
	var results []Result
	for _, s := range m.srcs {
		res, err := s.Search(ctx, q)
		if err == nil {
			results = append(results, res...)
		}
	}
	if len(results) == 0 {
		return results, err
	}
	return results, nil
}

func (m MultiSource) Download(ctx context.Context, r *Result) error {
	for _, s := range m.srcs {
		if s.IsSource(r.Source) {
			return s.Download(ctx, r)
		}
	}
	return fmt.Errorf("No match source %q", r.Source)
}

func (m MultiSource) IsSource(key SourceKey) bool {
	for _, s := range m.srcs {
		if s.IsSource(key) {
			return true
		}
	}
	return false
}

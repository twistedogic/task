package sourcegraph

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/twistedogic/task/pkg/search"
)

const (
	BASE_URL = "https://sourcegraph.com/search/stream"
)

type Match struct {
	Number int     `json:"lineNumber"`
	Offset [][]int `json:"offsetAndLengths"`
}

type Result struct {
	Repo    string  `json:"repository"`
	Commit  string  `json:"commit"`
	File    string  `json:"path"`
	Matches []Match `json:"lineMatches"`
}

func (r Result) toSearch() search.Result {
	var matches []search.Match
	for _, m := range r.Matches {
		for _, o := range m.Offset {
			matches = append(matches, search.Match{
				Line:   m.Number,
				Offset: o[0],
				Length: o[1],
			})
		}
	}
	res := search.Result{
		Repo:    r.Repo,
		Commit:  r.Commit,
		File:    r.File,
		Matches: matches,
	}
	if strings.HasPrefix(r.Repo, "github.com/") {
		repo := strings.TrimPrefix(r.Repo, "github.com/")
		path := filepath.Join(repo, r.Commit, r.File)
		res.Link = fmt.Sprintf("https://raw.githubusercontent.com/%s", path)
	}
	return res
}

func toQueryString(q search.Query) string {
	tokens := []string{"context:global"}
	if q.Repo != "" {
		tokens = append(tokens, fmt.Sprintf("repo:%s", q.Repo))
	}
	if q.File != "" {
		tokens = append(tokens, fmt.Sprintf("file:%s", q.File))
	}
	tokens = append(tokens, q.Term, "case:yes")
	return strings.Join(tokens, " ")
}

type Client struct {
	*http.Client
	BaseURL string
}

func New(base string) Client {
	return Client{Client: new(http.Client), BaseURL: base}
}

func NewWithDefault() Client {
	return New(BASE_URL)
}

func (c Client) toURL(q search.Query) string {
	u, _ := url.Parse(c.BaseURL)
	qs := u.Query()
	qs.Add("q", toQueryString(q))
	qs.Add("v", "V2")
	qs.Add("t", "literal")
	u.RawQuery = qs.Encode()
	return u.String()
}

func (c Client) getWithContext(ctx context.Context, u string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got %v from %s", res.StatusCode, u)
	}
	return res.Body, nil
}

func (c Client) fetchFile(ctx context.Context, r search.Result) ([]byte, error) {
	if r.Link == "" {
		return nil, nil
	}
	res, err := c.getWithContext(ctx, r.Link)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	return ioutil.ReadAll(res)
}

func (c Client) parseResults(ctx context.Context, r io.Reader) ([]search.Result, error) {
	results := make([]search.Result, 0)
	scanner := bufio.NewScanner(r)
	isMatch := false
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			t := scanner.Text()
			if t == "event: matches" {
				isMatch = true
				continue
			}
			if !isMatch {
				continue
			}
			isMatch = false
			matches := make([]Result, 0)
			if err := json.Unmarshal([]byte(strings.TrimLeft(t, "data: ")), &matches); err != nil {
				return nil, err
			}
			for _, m := range matches {
				results = append(results, m.toSearch())
			}
		}
	}
	g, gctx := errgroup.WithContext(ctx)
	for i, r := range results {
		i, r := i, r
		g.Go(func() error {
			b, err := c.fetchFile(gctx, r)
			if err == nil {
				results[i].Content = b
			}
			return err
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return results, nil
}

func (c Client) Search(ctx context.Context, q search.Query) ([]search.Result, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.toURL(q), nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return c.parseResults(ctx, res.Body)
}

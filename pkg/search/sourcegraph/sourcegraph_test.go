package sourcegraph

import (
	"context"
	"testing"

	"github.com/twistedogic/task/pkg/search"
)

func Test_Client(t *testing.T) {
	t.Skip()
	client := NewWithDefault()
	q := search.Query{
		Repo: "prometheus",
		File: "ast.go",
		Term: "Expr",
	}
	res, err := client.Search(context.TODO(), q)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}

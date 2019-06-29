package database

import (
	"io/ioutil"
	"os"
	"testing"
)

func Setup(t *testing.T) (string, func() error) {
	t.Helper()
	file, err := ioutil.TempFile("", "database_test")
	if err != nil {
		t.Fatal(err)
	}
	c := func() error {
		return os.Remove(file.Name())
	}
	return file.Name(), c
}

func TestDatabase(t *testing.T) {
	type item struct {
		bucket []byte
		key    []byte
		value  []byte
	}
	cases := map[string]struct {
		items []item
	}{
		"base": {
			items: []item{
				{[]byte("b1"), []byte("key1"), []byte(`{"test":1}`)},
				{[]byte("b1"), []byte("key2"), []byte(`{"test":2}`)},
				{[]byte("b1"), []byte("key3"), []byte(`{"test":3}`)},
				{[]byte("b1"), []byte("key4"), []byte(`{"test":4}`)},
			},
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			path, clean := Setup(t)
			defer clean()
			db, err := New(path)
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			newBucket := []byte(name)
			for _, item := range tc.items {
				if err := db.Set(item.bucket, item.key, item.value); err != nil {
					t.Fatal(err)
				}
			}
			for _, item := range tc.items {
				b, err := db.Get(item.bucket, item.key)
				if err != nil {
					t.Fatal(err)
				}
				if string(b) != string(item.value) {
					t.Fatalf("want %s got %s", item.value, b)
				}
			}
			for _, item := range tc.items {
				if err := db.Move(item.bucket, newBucket, item.key); err != nil {
					t.Fatal(err)
				}
			}
			allKeys, err := db.List(newBucket)
			if err != nil {
				t.Fatal(err)
			}
			if len(allKeys) != len(tc.items) {
				t.Fatalf("want %v got %v", tc.items, allKeys)
			}
		})
	}
}

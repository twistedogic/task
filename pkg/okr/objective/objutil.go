package objective

import (
	"encoding/json"

	"github.com/twistedogic/task/pkg/store"
)

type ByTargetDate []*Objective

func (t ByTargetDate) Len() int {
	return len(t)
}

func (t ByTargetDate) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t ByTargetDate) Less(i, j int) bool {
	return t[i].Target.Before(t[j].Target)
}

func List(db store.Store) ([]*Objective, error) {
	items, err := db.List([]byte(BucketName))
	if err != nil {
		return nil, err
	}
	objs := make([]*Objective, 0, len(items))
	for _, item := range items {
		o, err := Load(item.Value)
		if err != nil {
			return nil, err
		}
		objs = append(objs, o)
	}
	return objs, nil
}

func Load(b []byte) (*Objective, error) {
	o := new(Objective)
	if err := json.Unmarshal(b, o); err != nil {
		return o, err
	}
	o = o.LoadKeyResults()
	return o, nil
}

package okr

import (
	"sort"

	"github.com/manifoldco/promptui"
	"github.com/twistedogic/store"
)

type Store struct {
	db    store.Store
	index map[string]*obj.Objective
}

func New(s store.Store) (*Store, error) {
	objs, err := obj.List(s)
	if err != nil {
		return nil, err
	}
	index := make(map[string]*obj.Objective)
	for _, o := range objs {
		index[o.Name] = o
	}
	return &Store{db: s, index: index}, nil
}

func (s *Store) List() []*obj.Objective {
	objs := make([]*obj.Objective, 0, len(s.index))
	for _, v := range s.index {
		objs = append(objs, v)
	}
	sort.Sort(obj.ByTargetDate(objs))
	return objs
}

func (s *Store) Add() error {
	o, err := obj.New()
	if err != nil {
		return err
	}
	return o.Save(s.db)
}

func (s *Store) promptObjectives() (*obj.Objective, error) {
	keys := make([]string, 0, len(s.index))
	for k := range s.index {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	selectObj := promptui.Select{
		Label: "Objectives:",
		Items: keys,
	}
	_, key, err := selectObj.Run()
	if err != nil {
		return nil, err
	}
	return s.index[key], nil
}
func (s *Store) Edit() error {
	o, err := s.promptObjectives()
	if err != nil {
		return err
	}
	return o.Update(s.db)
}

func (s *Store) Update() error {
	o, err := s.promptObjectives()
	if err != nil {
		return err
	}
	if err := o.UpdateKR(); err != nil {
		return err
	}
	return o.Save(s.db)
}

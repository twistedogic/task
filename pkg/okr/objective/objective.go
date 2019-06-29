package objective

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/twistedogic/task/pkg/okr/keyresult"
	"github.com/twistedogic/task/pkg/prompt"
	"github.com/twistedogic/task/pkg/store"
)

const (
	format          = "# Objective@%s (%s)\n- %s\n## Key Results\n%s\n"
	timestampFormat = "2006-Jan-02"
	BucketName      = "objective"
)

type Objective struct {
	Name       string                          `json:"name"`
	Target     time.Time                       `json:"target"`
	Done       bool                            `json:"done"`
	Score      float64                         `json:"score"`
	KeyResults []*keyresult.KeyResult          `json:"key_results"`
	Index      map[string]*keyresult.KeyResult `json:"-"`
	Board      string                          `json:"board"`
	LastUpdate time.Time                       `json:"last_update"`
}

func New() (*Objective, error) {
	o := new(Objective)
	err := o.Prompt()
	return o, err
}

func (o *Objective) Prompt() error {
	board, err := prompt.StringWithDefault("Board", o.Board)
	if err != nil {
		return err
	}
	name, err := prompt.StringWithDefault("Name", o.Name)
	if err != nil {
		return err
	}
	target, err := prompt.Date("Target Date", timestampFormat)
	if err != nil {
		return err
	}
	o.Board, o.Name, o.Target = board, name, target
	for {
		if len(o.Index) != 0 {
			edit, err := prompt.YN("Edit KR")
			if err != nil {
				return err
			}
			if edit {
				if err := o.EditKR(); err != nil {
					return err
				}
				o = o.LoadKeyResults()
			}
		}
		add, err := prompt.YN("Add KR")
		if err != nil {
			return err
		}
		if !add {
			break
		}

		kr, err := keyresult.New()
		if err != nil {
			return err
		}
		o.KeyResults = append(o.KeyResults, kr)
		o = o.LoadKeyResults()
	}
	o.LastUpdate = time.Now()
	return nil
}

func (o *Objective) LoadKeyResults() *Objective {
	o.Index = make(map[string]*keyresult.KeyResult)
	for _, kr := range o.KeyResults {
		o.Index[kr.Name] = kr
	}
	keyResults := make([]*keyresult.KeyResult, 0, len(o.Index))
	for _, k := range o.Index {
		keyResults = append(keyResults, k)
	}
	o.KeyResults = keyResults
	return o
}

func (o *Objective) UpdateScore() *Objective {
	score := 0.0
	for _, kr := range o.KeyResults {
		score += kr.GetScore()
	}
	o.Score = score
	return o
}

func (o *Objective) String() string {
	target := o.Target.Format(timestampFormat)
	krList := "No Key Result defined"
	krs := make([]string, len(o.KeyResults))
	if len(o.KeyResults) > 0 {
		for i, kr := range o.KeyResults {
			krs[i] = fmt.Sprintf("%s", kr)
		}
		sort.Strings(krs)
		krList = strings.Join(krs, "\n")
	}
	return fmt.Sprintf(format, target, o.Board, o.Name, krList)
}

func (o *Objective) Key() []byte {
	return []byte(o.Name)
}

func (o *Objective) promptKRs() (*keyresult.KeyResult, error) {
	krList := make([]string, 0, len(o.Index))
	for k := range o.Index {
		krList = append(krList, k)
	}
	prompt := promptui.Select{Label: "Update Key Results", Items: krList}
	_, result, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	return o.Index[result], nil
}

func (o *Objective) UpdateKR() error {
	kr, err := o.promptKRs()
	if err != nil {
		return err
	}
	if err := kr.Update(); err != nil {
		return err
	}
	o = o.LoadKeyResults()
	return nil
}

func (o *Objective) EditKR() error {
	kr, err := o.promptKRs()
	if err != nil {
		return err
	}
	if err := kr.Prompt(); err != nil {
		return err
	}
	o = o.LoadKeyResults()
	return nil
}

func (o *Objective) Save(s store.Store) error {
	value, err := json.Marshal(o.LoadKeyResults())
	if err != nil {
		return err
	}
	return s.Set([]byte(BucketName), o.Key(), value)
}

func (o *Objective) Update(s store.Store) error {
	oldkey := o.Key()
	if err := o.Prompt(); err != nil {
		return err
	}
	if err := s.Delete([]byte(BucketName), oldkey); err != nil {
		return err
	}
	return o.Save(s)
}

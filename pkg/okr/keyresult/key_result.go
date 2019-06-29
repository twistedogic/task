package keyresult

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/twistedogic/task/pkg/prompt"
)

type KeyResult struct {
	Name       string    `json:"name"`
	LastUpdate time.Time `json:"last_update"`
	Total      int       `json:"total"`
	Completed  int       `json:"completed"`
	Done       bool      `json:"done"`
}

func Load(b []byte) (*KeyResult, error) {
	var kr *KeyResult
	err := json.Unmarshal(b, kr)
	return kr, err
}

func New() (*KeyResult, error) {
	kr := new(KeyResult)
	err := kr.Prompt()
	return kr, err
}

func (kr *KeyResult) GetScore() float64 {
	return float64(kr.Completed) / float64(kr.Total)
}

func (kr *KeyResult) IsDone() bool {
	kr.Done = kr.Completed >= kr.Total
	return kr.Done
}

func (kr *KeyResult) String() string {
	return fmt.Sprintf("- %d/%d %s", kr.Completed, kr.Total, kr.Name)
}

func (kr *KeyResult) Prompt() error {
	name, err := prompt.StringWithDefault("Key Result", kr.Name)
	if err != nil {
		return err
	}
	total, err := prompt.Int("Total")
	if err != nil {
		return err
	}
	kr.Name, kr.Total = name, total
	kr.LastUpdate = time.Now()
	return nil
}

func (kr *KeyResult) Update() error {
	complete, err := prompt.IntWithDefault(kr.Name, kr.Completed)
	if err != nil {
		return err
	}
	kr.Completed = complete
	return nil
}

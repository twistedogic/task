package okr

import (
	"fmt"
	"time"
)

type KeyResult struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	LastUpdate time.Time `json:"last_update"`
	Total      int       `json:"total"`
	Completed  int       `json:"completed"`
}

func (k *KeyResult) FilterValue() string { return k.Name }
func (k *KeyResult) Title() string       { return k.Name }
func (k *KeyResult) Description() string { return fmt.Sprintf("%d / %d", k.Completed, k.Total) }

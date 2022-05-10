package okr

import (
	"fmt"
	"time"
)

type Objective struct {
	Name       string
	LastUpdate time.Time
	KeyResults []*KeyResult
}

func (o *Objective) progress() float64 {
	completed, total := 0, 0
	for _, kr := range o.KeyResults {
		completed += kr.Completed
		total += kr.Total
	}
	ratio := float64(completed) / float64(total)
	return ratio * 100
}

func (o *Objective) FilterValue() string { return o.Name }
func (o *Objective) Title() string       { return o.Name }
func (o *Objective) Description() string { return fmt.Sprintf("%.f", o.progress()) }

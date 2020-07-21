package ev

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

func GetBreakeven(args []string) (float64, error) {
	if len(args) != 2 {
		return 0, fmt.Errorf("expect 2 inputs")
	}
	potString, pString := args[0], args[1]
	pot, err := strconv.ParseFloat(potString, 64)
	if err != nil {
		return 0, err
	}
	p, err := strconv.ParseFloat(pString, 64)
	if err != nil {
		return 0, err
	}
	if pot < p {
		pot, p = p, pot
	}
	return pot * p / (1.0 - p), nil
}

var RunCmd = &cobra.Command{
	Use:   "ev",
	Short: "expected breakeven",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := GetBreakeven(args)
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(out)
		}
	},
}

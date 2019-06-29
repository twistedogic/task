package prompt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/manifoldco/promptui"
)

func validateDate(format string) func(string) error {
	return func(input string) error {
		_, err := time.Parse(format, input)
		return err
	}
}

func validateYN(input string) error {
	if input == "y" || input == "n" {
		return nil
	}
	return fmt.Errorf("enter y/n")
}

func validateInt(input string) error {
	_, err := strconv.Atoi(input)
	return err
}

func Date(label, format string) (time.Time, error) {
	validate := validateDate(format)
	prompt := promptui.Prompt{
		Label:     "Target Date",
		Validate:  validate,
		AllowEdit: true,
	}
	res, _ := prompt.Run()
	return time.Parse(format, res)
}

func StringWithDefault(label, defaultInput string) (string, error) {
	prompt := promptui.Prompt{
		Label:     label,
		Default:   defaultInput,
		AllowEdit: true,
	}
	return prompt.Run()
}

func Int(label string) (int, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validateInt,
	}
	res, err := prompt.Run()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

func IntWithDefault(label string, i int) (int, error) {
	prompt := promptui.Prompt{
		Label:     label,
		Validate:  validateInt,
		Default:   strconv.Itoa(i),
		AllowEdit: true,
	}
	res, err := prompt.Run()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}
func YN(label string) (bool, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Default:  "y",
		Validate: validateYN,
	}
	res, err := prompt.Run()
	return res == "y", err
}

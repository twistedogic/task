package cheat

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

const BASE = "http://cht.sh/"

func GetAnswer(lang string, page int, quiet bool, query ...string) (string, error) {
	client := &http.Client{}
	qs := strings.Join(query, "+")
	u := fmt.Sprintf("%s/%s", BASE, qs)
	if lang != "" {
		u = fmt.Sprintf("%s/%s/%s", BASE, lang, qs)
	}
	if page > 0 {
		u = fmt.Sprintf("%s/%d", u, page)
	}
	if quiet {
		u = fmt.Sprintf("%s?Q", u)
	}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "curl")
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	return string(b), err
}

var langFlag string
var pageFlag int
var quietFlag bool

var RunCmd = &cobra.Command{
	Use:   "cht",
	Short: "cheat.sh",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := GetAnswer(langFlag, pageFlag, quietFlag, args...)
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(out)
		}
	},
}

func init() {
	RunCmd.Flags().StringVarP(&langFlag, "lang", "l", "", "target language")
	RunCmd.Flags().IntVarP(&pageFlag, "page", "p", 0, "page number")
	RunCmd.Flags().BoolVarP(&quietFlag, "quiet", "q", true, "quiet mode (code only)")
}

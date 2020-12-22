package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"testing"

	"github.com/manifoldco/promptui"
)

func TestCompareAndSync(t *testing.T) {
	s := map[string]interface{}{
		"default": "http://localhost:8080/",
	}
	viper.Set("remoteDefault", "default")
	viper.Set("remote", s)
	CompareAndSync("gg-example", "")
}

func TestPrompt(t *testing.T) {
	items := []string{"Vim", "Emacs", "Sublime", "VSCode", "Atom"}
	index := -1
	var result string
	var err error

	for index < 0 {
		prompt := promptui.SelectWithAdd{
			Label:    "What's your text editor",
			Items:    items,
			AddLabel: "Other",
		}

		index, result, err = prompt.Run()

		if index == -1 {
			items = append(items, result)
		}
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %s\n", result)
}

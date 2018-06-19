package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/manifoldco/promptui"
)

func main() {
	prompt := promptui.Prompt{
		Label: "input dept name's pinyin code",
		Validate: func(input string) error {
			if len(input) < 4 {
				return errors.New("length of input >= 4")
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(result)

	var depts [][3]string

	depts = [][3]string{
		{"2301", "h04", "妇产科病区"},
		{"2801", "h09", "外科病区"},
	}

	prompt2 := promptui.Select{
		Label: "select a dept",
		Items: depts,
	}

	index, result, err := prompt2.Run()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(index, result)
}

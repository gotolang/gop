package main

import (
	"fmt"

	"gopkg.in/AlecAivazis/survey.v1"
)

var qs = []*survey.Question{
	{
		Name:      "name",
		Prompt:    &survey.Input{Message: "what is your name?"},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
	{
		Name: "color",
		Prompt: &survey.Select{
			Message: "Choose a color:",
			Options: []string{"red", "blue", "green"},
			Default: "red",
		},
	},
	{
		Name:   "age",
		Prompt: &survey.Input{Message: "How old are you?"},
	},
}

func main() {
	answers := struct {
		Name          string
		FavoriteColor string `survey:"color"`
		Age           int
	}{} //匿名结构体初始化

	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%s chose %s.", answers.Name, answers.FavoriteColor)

}

package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/thoas/go-funk"
)

/*
{
  "key": "author",
  "desc": "作者名称",
  "type": "input"
  "conditions": ["a","b"]
},
*/

type Prompt struct {
	Questions []*survey.Question
	Answer    map[string]interface{}
	resource  []map[string]interface{}
}

func NewPrompt(resource []map[string]interface{}) Prompt {
	return Prompt{resource: resource}
}

func (p *Prompt) Parse() {
	var questions []*survey.Question
	for _, r := range p.resource {
		name := r["key"].(string)
		desc := r["desc"].(string)
		tp := r["type"].(string)
		item := survey.Question{
			Name:     name,
			Validate: survey.Required,
		}

		switch tp {
		case "radio":
			ops := r["conditions"].([]interface{})
			item.Prompt = &survey.Select{
				Message: desc,
				Options: toSlice(ops),
			}
		case "input":
			item.Prompt = &survey.Input{Message: desc}

		case "checkbox":
			ops := r["conditions"].([]interface{})
			item.Prompt = &survey.MultiSelect{
				Message: desc,
				Options: toSlice(ops),
			}
		default:

		}
		questions = append(questions, &item)
	}
	p.Questions = questions
}

func toSlice(o []interface{}) []string {
	return funk.Map(o, func(x interface{}) string {
		return x.(string)
	}).([]string)
}

func (p *Prompt) Run() error {
	answers := make(map[string]interface{})
	p.Answer = answers
	return survey.Ask(p.Questions, &answers)
}

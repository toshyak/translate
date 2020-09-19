package main

import (
	"encoding/json"
	"fmt"
)

// https://www.alfredapp.com/help/workflows/inputs/script-filter/json/
type alfredOutput struct {
	Items []alfredOutputItem `json:"items"`
}

type alfredOutputItem struct {
	Title    string            `json:"title"`
	Subtitle string            `json:"subtitle"`
	Arg      string            `json:"arg"`
	Icon     map[string]string `json:"icon"`
}

var translationServiceIcons = map[string]string{
	"aws":     "aws_translate.png",
	"speller": "speller_logo.png",
}

func newOutput() *alfredOutput {
	var t alfredOutput
	return &t
}

func (t *alfredOutput) add(text string, subtitle string, translationService string) {
	icon := map[string]string{"path": translationServiceIcons[translationService]}
	item := alfredOutputItem{Title: text, Arg: text, Subtitle: subtitle, Icon: icon}
	t.Items = append(t.Items, item)
	outputJSON, _ := json.Marshal(t)
	fmt.Print(string(outputJSON))
}

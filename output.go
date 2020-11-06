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
	Title        string            `json:"title"`
	Subtitle     string            `json:"subtitle"`
	Arg          string            `json:"arg"`
	Valid        bool              `json:"valid"`
	Autocomplete string            `json:"autocomplete"`
	Icon         map[string]string `json:"icon"`
}

var translationServiceIcons = map[string]string{
	"aws":     "aws_translate.png",
	"speller": "speller_logo.png",
	"ydict":   "ydict_logo.png",
}

func newOutput() *alfredOutput {
	var t alfredOutput
	return &t
}

func (t *alfredOutput) print() {
	outputJSON, _ := json.Marshal(t)
	fmt.Println(string(outputJSON))
}

func (t *alfredOutput) add(text string, subtitle string, translationService string, isValid bool) {
	icon := map[string]string{"path": translationServiceIcons[translationService]}
	item := alfredOutputItem{
		Title:        text,
		Arg:          text,
		Valid:        isValid,
		Autocomplete: text,
		Subtitle:     subtitle,
		Icon:         icon}
	t.Items = append(t.Items, item)
}

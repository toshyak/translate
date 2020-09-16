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
	Icon     map[string]string `json:"icon"`
}

func output(text string) {
	icon := map[string]string{"path": "aws_translate.png"}
	item := alfredOutputItem{Title: text, Icon: icon}
	output := alfredOutput{Items: []alfredOutputItem{item}}
	outputJSON, _ := json.Marshal(output)
	fmt.Print(string(outputJSON))
}

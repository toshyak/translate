package main

import (
	"log"
	"os"
	"strings"
	"toshyak/translate/aws"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Pass translated text as an argument")
	}
	textToTranslate := strings.Join(os.Args[1:], " ")
	output(aws.Translate(textToTranslate))
}

package main

import (
	"log"
	"os"
	"strings"
	"toshyak/translate/aws"
	"toshyak/translate/spelling"
	"unicode"
)

var translationDirections = map[string]string{
	"ru": "en",
	"en": "ru",
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Pass translated text as an argument")
	}
	textToTranslate := strings.Join(os.Args[1:], " ")
	sourceLanguage := getSourceLanguage(textToTranslate)
	spellingSuggestions, err := spelling.CheckSpelling(textToTranslate, sourceLanguage)
	if err != nil {
		log.Println("Failed to check spelling", err)
	}
	translatedText := aws.Translate(textToTranslate, sourceLanguage, translationDirections[sourceLanguage])
	out := newOutput()
	for _, s := range spellingSuggestions {
		out.add(s, "", "speller")
	}
	out.add(translatedText, "", "aws")
	out.print()
}

func getSourceLanguage(text string) string {
	f := func(r rune) bool {
		return unicode.Is(unicode.Cyrillic, r)
	}
	if strings.IndexFunc(text, f) != -1 {
		return "ru"
	}
	return "en"
}

package main

import (
	"fmt"

	"github.com/c-bata/go-prompt"

	customPrompt "github.com/sergalkin/gophkeeper/internal/client/prompt"
)

var (
	buildVersion string
	buildDate    string
)

func main() {
	fmt.Printf("Build version: %s\nBuild date: %s\n", buildVersion, buildDate)

	p := prompt.New(
		customPrompt.NewExecutor().Execute,
		customPrompt.NewCompleter().Complete,
		prompt.OptionTitle("Gophkeeper"),
		prompt.OptionPrefix(">>>"),
		prompt.OptionInputTextColor(prompt.Yellow),
	)
	p.Run()
}

package prompt

import "github.com/c-bata/go-prompt"

type Completer struct {
}

func NewCompleter() *Completer {
	return &Completer{}
}

// Complete - a list of suggestions.
func (c *Completer) Complete(d prompt.Document) []prompt.Suggest {
	var s []prompt.Suggest

	if d.FindStartOfPreviousWord() == 0 {
		s = []prompt.Suggest{
			{Text: "login", Description: "Authenticate user"},
			{Text: "logout", Description: "Logout authenticated user"},
			{Text: "register", Description: "Register new user"},
			{Text: "delete-user", Description: "Delete logged user"},
			{Text: "types", Description: "Get list of secret types available to be stored"},
			{Text: "create-auth", Description: "Create new login/pass secret"},
			{Text: "create-text", Description: "Create new text secret"},
			{Text: "create-binary", Description: "Create new binary secret"},
			{Text: "create-card", Description: "Create new card secret"},
			{Text: "get-secret", Description: "Retrieve stored secret"},
			{Text: "get-secret-binary", Description: "Retrieve stored binary secret"},
			{Text: "delete-secret", Description: "Retrieve stored secret"},
			{Text: "edit-secret", Description: "Edit stored secret"},
			{Text: "get-secrets-by-type", Description: "Retrieves list of secretes by their type"},
			{Text: "exit", Description: "Exit program"},
		}
	}

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

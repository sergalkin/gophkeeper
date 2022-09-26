package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sergalkin/gophkeeper/internal/client/app"
	"github.com/sergalkin/gophkeeper/internal/client/model"
)

var (
	clientApp *app.App

	buildVersion string
	buildDate    string
)

func main() {
	fmt.Printf("Build version: %s\nBuild date: %s\n", buildVersion, buildDate)

	appLoc, err := app.NewApp()
	if err != nil {
		fmt.Println(err)
		return
	}
	clientApp = appLoc

	p := prompt.New(
		executor, completer,
		prompt.OptionTitle("Gophkeeper"),
		prompt.OptionPrefix(">>>"),
		prompt.OptionInputTextColor(prompt.Yellow),
	)
	p.Run()
}

// executor - executes proper function based on entered text in terminal.
func executor(s string) {
	var isForce bool

	setCommand, options := getCommandArgsAndOptions(s)
	if options["force"] || options["f"] {
		isForce = true
	}

	switch setCommand[0] {
	case "login":
		switch len(setCommand) - 1 {
		case 1:
			fmt.Println("you have to write password")
			return
		case 0:
			fmt.Println("you have to write login and password")
			return
		}

		user := model.User{Login: setCommand[1], Password: setCommand[2]}
		if err := clientApp.UserService.Login(user); err != nil {
			switch status.Code(err) {
			case codes.NotFound:
				fmt.Println("User not found")
				return
			default:
				fmt.Println(err)
				return
			}
		}

		// firstly we sync all on start up
		clientApp.Syncer.SyncAll()

		// then we spawn goroutin with cron job to sync data every minute
		go clientApp.Cron.Run()

		fmt.Println("successfully authorized")
		return
	case "register":
		switch len(setCommand) - 1 {
		case 1:
			fmt.Println("you have to write password")
			return
		case 0:
			fmt.Println("you have to write login and password")
			return
		}
		user := model.User{Login: setCommand[1], Password: setCommand[2]}
		if err := clientApp.UserService.Register(user); err != nil {
			switch status.Code(err) {
			case codes.InvalidArgument:
				fmt.Println("Data validation error or user already exists")
				return
			default:
				fmt.Println(err)
				return
			}
		}

		fmt.Println("User successfully created. You are logged in.")

		return
	case "delete-user":
		if err := clientApp.UserService.Delete(); err != nil {
			fmt.Println(err)
		}
		return
	case "logout":
		clientApp.UserService.Logout()

		clientApp.Cron.Stop()

		clientApp.Storage.ResetStorage()

		fmt.Println("You successfully logged out")
		return
	case "secret-types":
		secrets, err := clientApp.SecretTypeService.List()
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, secret := range secrets.Secrets {
			fmt.Printf("%+v\n", secret)
		}

		return
	case "create-auth-secret":
		switch len(setCommand) - 1 {
		case 2:
			fmt.Println("you have to write password")
			return
		case 1:
			fmt.Println("you have to write login, password")
			return
		case 0:
			fmt.Println("you have to write title, login, password")
			return
		}

		m := model.LoginPassSecret{
			Title:      setCommand[1],
			RecordType: 1,
			Login:      setCommand[2],
			Password:   setCommand[3],
		}

		cont, errMarshal := json.Marshal(m)
		if errMarshal != nil {
			fmt.Println(errMarshal)
			return
		}

		if err := clientApp.SecretService.CreateSecret(m.Title, 1, string(cont)); err != nil {
			fmt.Println(err)
		}

		return
	case "create-text-secret":
		switch len(setCommand) - 1 {
		case 1:
			fmt.Println("you have to write text")
			return
		case 0:
			fmt.Println("you have to write title and text")
			return
		}
		m := model.TextSecret{
			Title:      setCommand[1],
			RecordType: 2,
			Text:       strings.Join(setCommand[2:], " "),
		}

		marshal, errMarshal := json.Marshal(m)
		if errMarshal != nil {
			fmt.Println(errMarshal)
			return
		}

		if err := clientApp.SecretService.CreateSecret(m.Title, 2, string(marshal)); err != nil {
			fmt.Println(err)
		}

		return
	case "create-binary-secret":
		switch len(setCommand) - 1 {
		case 1:
			fmt.Println("you have to write file path")
			return
		case 0:
			fmt.Println("you have to write title and file path")
			return
		}

		m := model.FileSecret{
			Title:      setCommand[1],
			RecordType: 3,
			Path:       setCommand[2],
		}

		f, err := os.Open(m.Path)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		data, errData := ioutil.ReadAll(f)
		if errData != nil {
			fmt.Println(errData)
			return
		}

		errCreate := clientApp.SecretService.CreateSecret(m.Title, m.RecordType, string(data))
		if errCreate != nil {
			fmt.Println(errCreate)
		}

		return
	case "create-card-secret":
		switch len(setCommand) - 1 {
		case 3:
			fmt.Println("you have to write due date")
			return
		case 2:
			fmt.Println("you have to write cvv, due date")
			return
		case 1:
			fmt.Println("you have to write card number, cvv, due date")
			return
		case 0:
			fmt.Println("you have to write title, card number, cvv, due date")
			return
		}
		cardModel := model.CardSecret{
			Title:      setCommand[1],
			RecordType: 4,
			CardNumber: setCommand[2],
			CVV:        setCommand[3],
			Due:        setCommand[4],
		}
		cont, er := json.Marshal(cardModel)
		if er != nil {
			fmt.Println(er)
			return
		}

		if err := clientApp.SecretService.CreateSecret(cardModel.Title, cardModel.RecordType, string(cont)); err != nil {
			fmt.Println(err)
		}

		return
	case "get-secret":
		switch len(setCommand) - 1 {
		case 0:
			fmt.Println("you have to write ID of secret")
			return
		}

		id, convErr := strconv.Atoi(setCommand[1])
		if convErr != nil {
			fmt.Println(convErr)
			return
		}

		if _, err := clientApp.SecretService.GetSecret(id); err != nil {
			switch status.Code(err) {
			case codes.NotFound:
				fmt.Println("Secret not found")
				return
			default:
				fmt.Println(err)
				return
			}
		}

		return
	case "get-binary-secret":
		switch len(setCommand) - 1 {
		case 0:
			fmt.Println("you have to write ID and location")
			return
		case 1:
			fmt.Println("you have to write location")
			return
		}

		id, errConv := strconv.Atoi(setCommand[1])
		if errConv != nil {
			fmt.Println(errConv)
			return
		}

		err := clientApp.SecretService.GetBinarySecret(id, setCommand[2])
		if err != nil {
			fmt.Println(err)
			return
		}

		return
	case "delete-secret":
		switch len(setCommand) - 1 {
		case 0:
			fmt.Println("you have to write ID of secret")
			return
		}

		id, convErr := strconv.Atoi(setCommand[1])
		if convErr != nil {
			fmt.Println(convErr)
			return
		}

		if err := clientApp.SecretService.DeleteSecret(id); err != nil {
			fmt.Println(err)
		}

		return
	case "edit-secret":
		var (
			recordType int
			id         int
			converted  []byte
		)

		numArgs := len(setCommand) - 1
		if numArgs >= 3 {
			recordType, errConv := strconv.Atoi(setCommand[3])
			if errConv != nil {
				fmt.Println(errConv)
				return
			}

			id, errConv = strconv.Atoi(setCommand[1])
			if errConv != nil {
				fmt.Println(errConv)
				return
			}

			switch recordType {
			case 1:
				switch numArgs {
				case 4:
					fmt.Println("you have to write password")
					return
				case 3:
					fmt.Println("you have to write login,password")
					return
				default:
					converted, errConv = json.Marshal(model.LoginPassSecret{
						Id:         id,
						Title:      setCommand[2],
						RecordType: 1,
						Login:      setCommand[4],
						Password:   setCommand[5],
					})
					if errConv != nil {
						fmt.Println(errConv)
						return
					}
				}
			case 2:
				switch numArgs {
				case 3:
					fmt.Println("you have to write text")
					return
				default:
					converted, errConv = json.Marshal(model.TextSecret{
						Id:         id,
						Title:      setCommand[2],
						RecordType: 2,
						Text:       strings.Join(setCommand[4:], " "),
					})
					if errConv != nil {
						fmt.Println(errConv)
						return
					}
				}
			case 3:
				switch numArgs {
				case 3:
					fmt.Println("you have to write file")
					return
				default:
					converted, errConv = json.Marshal(model.FileSecret{
						Id:         id,
						Title:      setCommand[2],
						RecordType: 3,
						Path:       setCommand[3],
					})
					if errConv != nil {
						fmt.Println(errConv)
						return
					}
				}
			case 4:
				switch numArgs {
				case 5:
					fmt.Println("you have to write due")
					return
				case 4:
					fmt.Println("you have to write cvv,due")
					return
				case 3:
					fmt.Println("you have to write card number,cvv,due")
					return
				default:
					converted, errConv = json.Marshal(model.CardSecret{
						Id:         id,
						Title:      setCommand[2],
						RecordType: 4,
						CardNumber: setCommand[4],
						CVV:        setCommand[5],
						Due:        setCommand[6],
					})
					if errConv != nil {
						fmt.Println(errConv)
						return
					}
				}
			}
		} else {
			fmt.Println("you have to write ID of secret, title, type and record fields")
			return
		}

		if err := clientApp.SecretService.EditSecret(id, setCommand[3], recordType, string(converted), isForce); err != nil {
			st, _ := status.FromError(err)

			fmt.Println(st.Message())

			if st.Code() == codes.FailedPrecondition {
				fmt.Println("starting re-sync")

				clientApp.Syncer.SyncAll()

				fmt.Println("re-sync ended")
			}

			return
		}

		return
	case "get-secret-list-by-type":
		switch len(setCommand) - 1 {
		case 0:
			fmt.Println("you have to write type ID of secret")
			return
		}

		id, convErr := strconv.Atoi(setCommand[1])
		if convErr != nil {
			fmt.Println(convErr)
			return
		}

		list, err := clientApp.SecretService.GetListOfSecretes(id)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, secret := range list {
			fmt.Printf("ID:%v Title: %v\n", secret.Id, secret.Title)
		}

		return
	case "exit":
		fmt.Println("application is closing")

		clientApp.Cancel()
		clientApp.Cron.Stop()

		os.Exit(0)
	}
}

// completer - a list of suggestions.
func completer(d prompt.Document) []prompt.Suggest {
	var s []prompt.Suggest

	if d.FindStartOfPreviousWord() == 0 {
		s = []prompt.Suggest{
			{Text: "login", Description: "Authenticate user"},
			{Text: "logout", Description: "Logout authenticated user"},
			{Text: "register", Description: "Register new user"},
			{Text: "delete-user", Description: "Delete logged user"},
			{Text: "secret-types", Description: "Get list of secret types available to be stored"},
			{Text: "create-auth-secret", Description: "Create new login/pass secret"},
			{Text: "create-text-secret", Description: "Create new text secret"},
			{Text: "create-binary-secret", Description: "Create new binary secret"},
			{Text: "create-card-secret", Description: "Create new card secret"},
			{Text: "get-secret", Description: "Retrieve stored secret"},
			{Text: "get-binary-secret", Description: "Retrieve stored binary secret"},
			{Text: "delete-secret", Description: "Retrieve stored secret"},
			{Text: "edit-secret", Description: "Edit stored secret"},
			{Text: "get-secret-list-by-type", Description: "Retrieves list of secretes by their type"},
			{Text: "exit", Description: "Exit program"},
		}
	}

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func getCommandArgsAndOptions(s string) ([]string, map[string]bool) {
	s = strings.TrimSpace(s)

	setCommand := strings.Split(s, " ")

	l := len(setCommand)

	filtered := make([]string, 0, l)
	options := make(map[string]bool)

	for i := 0; i < len(setCommand); i++ {
		if strings.HasPrefix(setCommand[i], "-") {
			opt := strings.TrimPrefix(setCommand[i], "--")
			opt = strings.TrimPrefix(opt, "-")

			optSplited := strings.Split(opt, "=")
			options[optSplited[0]] = true

			continue
		}
		filtered = append(filtered, setCommand[i])
	}

	return filtered, options
}

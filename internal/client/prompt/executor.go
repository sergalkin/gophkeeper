package prompt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sergalkin/gophkeeper/internal/client/app"
	"github.com/sergalkin/gophkeeper/internal/client/model"
)

type Executor struct {
	app *app.App
}

func NewExecutor() *Executor {
	appL, err := app.NewApp()
	if err != nil {
		panic(err)
	}

	return &Executor{app: appL}
}

func (e *Executor) Execute(s string) {
	var isForce bool

	setCommand, options := getCommandArgsAndOptions(s)
	if options["force"] || options["f"] {
		isForce = true
	}

	switch setCommand[0] {
	case "login":
		if err := e.login(setCommand); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("successfully authorized")
		return
	case "register":
		if err := e.register(setCommand); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("User successfully created. You are logged in.")

		return
	case "delete-user":
		if err := e.deleteUser(); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("you successfully deleted account and logged out")

		return
	case "logout":
		if err := e.logout(); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("you successfully logged out")
		return
	case "types":
		types, err := e.types()
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, t := range types {
			fmt.Printf("%+v\n", t)
		}

		return
	case "create-auth":
		if err := e.createAuth(setCommand); err != nil {
			fmt.Println(err)
			return
		}
		return
	case "create-text":
		if err := e.createText(setCommand); err != nil {
			fmt.Println(err)
			return
		}

		return
	case "create-binary":
		if err := e.createBinary(setCommand); err != nil {
			fmt.Println(err)
			return
		}

		return
	case "create-card":
		if err := e.createCard(setCommand); err != nil {
			fmt.Println(err)
			return
		}

		return
	case "delete-secret":
		if err := e.deleteSecret(setCommand); err != nil {
			fmt.Println(err)
			return
		}

		return
	case "get-secrets-by-type":
		list, err := e.getSecretsByTypeId(setCommand)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, secret := range list {
			fmt.Printf("ID:%v Title: %v\n", secret.Id, secret.Title)
		}

		return
	case "get-secret":
		secret, err := e.getSecret(setCommand)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Content:%+v\n", secret)

		return
	case "get-secret-binary":
		if err := e.getSecretBinary(setCommand); err != nil {
			fmt.Println(err)
			return
		}
		return
	case "edit-secret":
		if err := e.editSecret(setCommand, isForce); err != nil {
			fmt.Println(err)
			return
		}

		return
	case "exit":
		fmt.Println("bye bye...application is closing")

		e.app.Cancel()
		e.app.Cron.Stop()

		os.Exit(0)
	}
}

// login - is executor for "login" case in Execute method.
func (e *Executor) login(args []string) error {
	switch len(args) - 1 {
	case 1:
		return fmt.Errorf("validation error: Password is missing")
	case 0:
		return fmt.Errorf("validation error: Login and Password is missing")
	}

	user := model.User{Login: args[1], Password: args[2]}

	if err := e.app.UserService.Login(user); err != nil {
		st, _ := status.FromError(err)

		switch st.Code() {
		case codes.NotFound:
			return fmt.Errorf("error: User not found")
		default:
			return fmt.Errorf("error:" + st.Message())
		}
	}

	// firstly we sync all on start up
	e.app.Syncer.SyncAll()

	// then we spawn goroutin with cron job to sync data every minute
	go e.app.Cron.Run()

	return nil
}

// register - is executor for "register" case in Execute method.
func (e *Executor) register(args []string) error {
	switch len(args) - 1 {
	case 1:
		return fmt.Errorf("validation error: Password is missing")
	case 0:
		return fmt.Errorf("validation error: Login and Password is missing")
	}

	user := model.User{Login: args[1], Password: args[2]}
	if err := e.app.UserService.Register(user); err != nil {
		switch status.Code(err) {
		case codes.InvalidArgument:
			return fmt.Errorf("error: data is invalid or user already exists")
		default:
			return err
		}
	}

	return nil
}

// deleteUser - is executor for "delete-user" case in Execute method.
func (e *Executor) deleteUser() error {
	return e.app.UserService.Delete()
}

// logout - is executor for "types" case in Execute method.
func (e *Executor) logout() error {
	e.app.UserService.Logout()

	e.app.Cron.Stop()

	e.app.Storage.ResetStorage()

	return nil
}

// types - is executor for "types" case in Execute method.
func (e *Executor) types() ([]model.SecretType, error) {
	secrets, err := e.app.SecretTypeService.List()
	if err != nil {
		return nil, err
	}

	var models []model.SecretType
	for _, secret := range secrets.Secrets {

		models = append(models, model.SecretType{
			Id:    int(secret.Id),
			Title: secret.Title,
		})
	}

	return models, nil
}

// createAuth - is executor for "create-auth" case in Execute method.
func (e *Executor) createAuth(args []string) error {
	switch len(args) - 1 {
	case 2:
		return fmt.Errorf("validation error: Password is missing")
	case 1:
		return fmt.Errorf("validation error: Login and Password is missing")
	case 0:
		return fmt.Errorf("validation error: Title, Login, Password is missing")
	}

	m := model.LoginPassSecret{
		Title:      args[1],
		RecordType: 1,
		Login:      args[2],
		Password:   args[3],
	}

	cont, errMarshal := json.Marshal(m)
	if errMarshal != nil {
		return errMarshal
	}

	if err := e.app.SecretService.CreateSecret(m.Title, 1, string(cont)); err != nil {
		return err
	}

	return nil
}

// createText - is executor for "create-text" case in Execute method.
func (e *Executor) createText(args []string) error {
	switch len(args) - 1 {
	case 1:
		return fmt.Errorf("validation error: Text is missing")
	case 0:
		return fmt.Errorf("validation error: Title and Text is missing")
	}

	m := model.TextSecret{
		Title:      args[1],
		RecordType: 2,
		Text:       strings.Join(args[2:], " "),
	}

	marshal, errMarshal := json.Marshal(m)
	if errMarshal != nil {
		return errMarshal
	}

	if err := e.app.SecretService.CreateSecret(m.Title, 2, string(marshal)); err != nil {
		return err
	}

	return nil
}

// createBinary - is executor for "create-binary" case in Execute method.
func (e *Executor) createBinary(args []string) error {
	switch len(args) - 1 {
	case 1:
		return fmt.Errorf("validation error: Filepath is missing")
	case 0:
		return fmt.Errorf("validation error: Title and Filepath is missing")
	}

	m := model.FileSecret{
		Title:      args[1],
		RecordType: 3,
		Path:       args[2],
	}

	f, err := os.Open(m.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	data, errData := ioutil.ReadAll(f)
	if errData != nil {
		return errData
	}

	errCreate := e.app.SecretService.CreateSecret(m.Title, m.RecordType, string(data))
	if errCreate != nil {
		return errCreate
	}

	return nil
}

// createCard - is executor for "create-card" case in Execute method.
func (e *Executor) createCard(args []string) error {
	switch len(args) - 1 {
	case 3:
		return fmt.Errorf("validation error: Due date is missing")
	case 2:
		return fmt.Errorf("validation error: CVV and Due date is missing")
	case 1:
		return fmt.Errorf("validation error: Card number, CVV, Due date is missing")
	case 0:
		return fmt.Errorf("validation error: Title, Card number, CVV, Due date is missing")
	}

	cardModel := model.CardSecret{
		Title:      args[1],
		RecordType: 4,
		CardNumber: args[2],
		CVV:        args[3],
		Due:        args[4],
	}

	cont, er := json.Marshal(cardModel)
	if er != nil {
		return er
	}

	if err := e.app.SecretService.CreateSecret(cardModel.Title, cardModel.RecordType, string(cont)); err != nil {
		return err
	}

	return nil
}

// deleteSecret - is executor for "delete-secret" case in Execute method.
func (e *Executor) deleteSecret(args []string) error {
	switch len(args) - 1 {
	case 0:
		return fmt.Errorf("validation error: Secret ID is missing")
	}

	id, convErr := strconv.Atoi(args[1])
	if convErr != nil {
		return convErr
	}

	if err := e.app.SecretService.DeleteSecret(id); err != nil {
		return err
	}

	return nil
}

// getSecretsByTypeId - is executor for "get-secrets-by-type" case in Execute method.
func (e *Executor) getSecretsByTypeId(args []string) ([]model.SecretList, error) {
	switch len(args) - 1 {
	case 0:
		return nil, fmt.Errorf("validation error: Secret Type ID is missing")
	}

	id, convErr := strconv.Atoi(args[1])
	if convErr != nil {
		return nil, convErr
	}

	list, err := e.app.SecretService.GetListOfSecretes(id)
	if err != nil {
		return nil, err
	}

	var models []model.SecretList
	for _, secret := range list {
		models = append(models, model.SecretList{
			Id:    int(secret.Id),
			Title: secret.Title,
		})
	}

	return models, nil
}

// getSecret - is executor for "get-secret" case in Execute method.
func (e *Executor) getSecret(args []string) (interface{}, error) {
	switch len(args) - 1 {
	case 0:
		return nil, fmt.Errorf("validation error: Secret ID is missing")
	}

	id, convErr := strconv.Atoi(args[1])
	if convErr != nil {
		return nil, convErr
	}

	secret, err := e.app.SecretService.GetSecret(id)
	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.NotFound:
			return nil, fmt.Errorf(st.Message())
		default:
			return nil, err
		}
	}

	return secret, nil
}

// getSecretBinary - is executor for "get-secret-binary" case in Execute method.
func (e *Executor) getSecretBinary(args []string) error {
	switch len(args) - 1 {
	case 0:
		return fmt.Errorf("validation error: Secret ID and Path is missing")
	case 1:
		return fmt.Errorf("validation error: Path is missing")
	}

	id, errConv := strconv.Atoi(args[1])
	if errConv != nil {
		return errConv
	}

	err := e.app.SecretService.GetBinarySecret(id, args[2])
	if err != nil {
		return err
	}

	return nil
}

// editSecret - is executor for "edit-secret" case in Execute method.
func (e *Executor) editSecret(args []string, isForce bool) error {
	var (
		recordType int
		id         int
		converted  []byte
	)

	numArgs := len(args) - 1
	if numArgs >= 3 {
		recordType, errConv := strconv.Atoi(args[3])
		if errConv != nil {
			return errConv
		}

		id, errConv = strconv.Atoi(args[1])
		if errConv != nil {
			return errConv
		}

		switch recordType {
		case 1:
			switch numArgs {
			case 4:
				return fmt.Errorf("validation error: Password is missing")
			case 3:
				return fmt.Errorf("validation error: Login and Password is missing")
			default:
				converted, errConv = json.Marshal(model.LoginPassSecret{
					Id:         id,
					Title:      args[2],
					RecordType: 1,
					Login:      args[4],
					Password:   args[5],
				})
				if errConv != nil {
					return errConv
				}
			}
		case 2:
			switch numArgs {
			case 3:
				return fmt.Errorf("validation error: Text is missing")
			default:
				converted, errConv = json.Marshal(model.TextSecret{
					Id:         id,
					Title:      args[2],
					RecordType: 2,
					Text:       strings.Join(args[4:], " "),
				})
				if errConv != nil {
					return errConv
				}
			}
		case 3:
			switch numArgs {
			case 3:
				return fmt.Errorf("validation error: Filepath is missing")
			default:
				converted, errConv = json.Marshal(model.FileSecret{
					Id:         id,
					Title:      args[2],
					RecordType: 3,
					Path:       args[3],
				})
				if errConv != nil {
					return errConv
				}
			}
		case 4:
			switch numArgs {
			case 5:
				return fmt.Errorf("validation error: Due date is missing")
			case 4:
				return fmt.Errorf("validation error: CVV and Due date is missing")
			case 3:
				return fmt.Errorf("validation error: Card number, CVV and Due date is missing")
			default:
				converted, errConv = json.Marshal(model.CardSecret{
					Id:         id,
					Title:      args[2],
					RecordType: 4,
					CardNumber: args[4],
					CVV:        args[5],
					Due:        args[6],
				})
				if errConv != nil {
					return errConv
				}
			}
		}
	} else {
		return fmt.Errorf("validation error: Secret ID, Title, Secret Type ID and secret fields is missing")
	}

	if err := e.app.SecretService.EditSecret(id, args[3], recordType, string(converted), isForce); err != nil {
		st, _ := status.FromError(err)

		fmt.Println(st.Message())

		if st.Code() == codes.FailedPrecondition {
			fmt.Println("starting re-sync")

			e.app.Syncer.SyncAll()

			fmt.Println("re-sync ended")
		}

		return nil
	}

	return nil
}

// getCommandArgsAndOptions - splits args to command args and options
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

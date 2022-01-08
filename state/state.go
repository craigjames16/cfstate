package state

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/craigjames16/cfstate/aws"
	"github.com/craigjames16/cfstate/utils"
	"github.com/urfave/cli/v2"
)

type App struct {
	Name                string
	Template            string
	Config              string
	StackID             string
	CurrentTemplateHash string
	CurrentConfigHash   string
}

type AppState struct {
	App    App
	Status status
}

type status string

var (
	NotCreated status = "NOT_CREATED"
	OK         status = "OK"
	Diff       status = "DIFF"
)

func getState() (apps []App, err error) {
	var (
		appsJson []byte
	)

	appsJson, err = os.ReadFile("./state.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(appsJson, &apps)
	if err != nil {
		return nil, err
	}

	return apps, nil
}

func checkState() (state []AppState, err error) {
	var (
		apps         []App
		templateHash string
		configHash   string
	)
	apps, err = getState()
	if err != nil {
		return nil, err
	}

	for _, app := range apps {
		fmt.Println(app.StackID, app.Name)
		if app.StackID == "" {
			state = append(state, AppState{App: app, Status: NotCreated})
			continue
		}

		templateHash, err = utils.GetFileHash(app.Template)
		if err != nil {
			return nil, err
		}

		configHash, err = utils.GetFileHash(app.Config)
		if err != nil {
			return nil, err
		}

		if templateHash == app.CurrentTemplateHash && configHash == app.CurrentConfigHash {
			state = append(state, AppState{App: app, Status: OK})
		} else {
			state = append(state, AppState{App: app, Status: Diff})
		}
	}

	return state, nil
}

func CheckState(c *cli.Context) error {
	var (
		err error
	)

	states, err := checkState()
	if err != nil {
		return err
	}

	for _, state := range states {
		fmt.Printf("%s: %s\n", state.App.Name, state.Status)
	}

	return nil
}

func SyncState(c *cli.Context) (err error) {
	var (
		opOutput aws.CreateUpdateOutput
	)

	states, err := checkState()
	if err != nil {
		return err
	}

	for _, state := range states {
		fmt.Println(state.App.Name, state.Status)
		switch state.Status {
		case NotCreated:
			opOutput, err = aws.CreateStack(aws.AppInput{
				Name:     state.App.Name,
				Template: state.App.Template,
				Config:   state.App.Config,
			})
		case OK:
			fmt.Println("OK")
			continue
		case Diff:
			opOutput, err = aws.UpdateStack(aws.AppInput{
				StackID:  state.App.StackID,
				Template: state.App.Template,
				Config:   state.App.Config,
			})
		}

		if err != nil {
			return err
		}

		utils.Must(updateState(state.App, opOutput))

	}

	return nil
}

func updateState(app App, opOutput aws.CreateUpdateOutput) (err error) {
	var (
		state        []App
		newState     []App
		templateHash string
		configHash   string
		newStateJSON []byte
	)

	now := time.Now()
	sec := now.Unix()

	state, err = getState()

	if err != nil {
		return err
	}

	for _, appState := range state {
		if appState.Name == app.Name {
			templateHash, err = utils.GetFileHash(app.Template)
			utils.Must(err)

			configHash, err = utils.GetFileHash(app.Config)
			utils.Must(err)

			newState = append(newState, App{
				Name:                appState.Name,
				Template:            appState.Template,
				Config:              appState.Config,
				StackID:             opOutput.StackID,
				CurrentTemplateHash: templateHash,
				CurrentConfigHash:   configHash,
			})

		} else {
			newState = append(newState, appState)
		}

	}

	utils.Must(os.Rename("./state.json", fmt.Sprintf("./prev_states/state-%d.json", sec)))

	newStateJSON, err = json.Marshal(newState)
	utils.Must(err)

	utils.Must(os.WriteFile("./state.json", newStateJSON, 0644))

	return err
}

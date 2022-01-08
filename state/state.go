package state

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/craigjames16/cfstate/aws"
	"github.com/craigjames16/cfstate/github"
	"github.com/craigjames16/cfstate/utils"
	"github.com/urfave/cli/v2"
)

var REPO_BASE string = "/tmp"

type App struct {
	Name                string
	Template            string
	Config              string
	StackID             string
	CurrentTemplateHash string
	CurrentConfigHash   string
}

type AppState struct {
	App              App
	TemplateLocation string
	ConfigLocation   string
	RepoURL          string
	Status           status
}

type Repo struct {
	RepoURL string
	Apps    []App
}

type status string

var (
	NotCreated status = "NOT_CREATED"
	OK         status = "OK"
	Diff       status = "DIFF"
)

func getState() (repos []Repo, err error) {
	var (
		appsJson []byte
	)

	appsJson, err = os.ReadFile("./state.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(appsJson, &repos)
	if err != nil {
		return nil, err
	}

	return repos, nil
}

func checkState() (appStates []AppState, err error) {
	var (
		state        []Repo
		templateHash string
		configHash   string
		repoLocation string
		appState     AppState
	)

	state, err = getState()
	if err != nil {
		return nil, err
	}

	for _, repo := range state {
		if repoLocation, err = github.GetRepo(repo.RepoURL); err != nil {
			return nil, err
		}

		for _, app := range repo.Apps {
			templateLocation := fmt.Sprintf("%s%s", repoLocation, app.Template)
			configLocation := fmt.Sprintf("%s%s", repoLocation, app.Config)

			appState = AppState{
				App:              app,
				TemplateLocation: templateLocation,
				ConfigLocation:   configLocation,
				RepoURL:          repo.RepoURL,
			}

			if app.StackID == "" {
				appState.Status = NotCreated
				appStates = append(appStates, appState)
				continue
			}

			templateHash, err = utils.GetFileHash(templateLocation)
			if err != nil {
				return nil, err
			}

			configHash, err = utils.GetFileHash(configLocation)
			if err != nil {
				return nil, err
			}

			if templateHash == app.CurrentTemplateHash && configHash == app.CurrentConfigHash {
				appState.Status = OK
			} else {
				appState.Status = Diff
			}

			appStates = append(appStates, appState)
		}
	}

	return appStates, nil
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
		switch state.Status {
		case NotCreated:
			opOutput, err = aws.CreateStack(aws.AppInput{
				Name:     state.App.Name,
				Template: state.TemplateLocation,
				Config:   state.ConfigLocation,
			})
		case OK:
			fmt.Println("OK")
			continue
		case Diff:
			opOutput, err = aws.UpdateStack(aws.AppInput{
				StackID:  state.App.StackID,
				Template: state.TemplateLocation,
				Config:   state.ConfigLocation,
			})
		}

		if err != nil {
			return err
		}

		utils.Must(updateState(state, opOutput))

	}

	return nil
}

func updateState(appState AppState, opOutput aws.CreateUpdateOutput) (err error) {
	var (
		state        []Repo
		newState     []Repo
		newAppState  []App
		templateHash string
		configHash   string
		newStateJSON []byte
	)

	now := time.Now()
	sec := now.Unix()

	if state, err = getState(); err != nil {
		return err
	}

	for _, repo := range state {
		if repo.RepoURL == appState.RepoURL {
			for _, app := range repo.Apps {
				if app.Name == appState.App.Name {
					templateHash, err = utils.GetFileHash(appState.TemplateLocation)
					utils.Must(err)

					configHash, err = utils.GetFileHash(appState.ConfigLocation)
					utils.Must(err)

					newAppState = append(newAppState, App{
						Name:                app.Name,
						Template:            app.Template,
						Config:              app.Config,
						StackID:             opOutput.StackID,
						CurrentTemplateHash: templateHash,
						CurrentConfigHash:   configHash,
					})

				} else {
					newAppState = append(newAppState, app)
				}
			}

			newState = append(newState, Repo{
				RepoURL: repo.RepoURL,
				Apps:    newAppState,
			})

		} else {
			newState = append(newState, repo)
		}

	}

	utils.Must(os.Rename("./state.json", fmt.Sprintf("./prev_states/state-%d.json", sec)))

	newStateJSON, err = json.Marshal(newState)
	utils.Must(err)

	utils.Must(os.WriteFile("./state.json", newStateJSON, 0644))

	return err
}

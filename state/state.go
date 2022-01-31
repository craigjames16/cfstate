package state

import (
	"encoding/json"
	"fmt"

	"github.com/craigjames16/cfstate/aws"
	"github.com/craigjames16/cfstate/backends"
	"github.com/craigjames16/cfstate/github"
	"github.com/craigjames16/cfstate/utils"
)

type App struct {
	Name                string
	Template            string
	Config              string
	StackID             string
	CurrentTemplateHash string
	CurrentConfigHash   string
}

type AppStatus struct {
	App              App
	TemplateLocation string
	ConfigLocation   string
	RepoURL          string
	Status           status
}

type Repo struct {
	RepoURL  string
	RepoName string
	Apps     []App
}

type State []Repo

type status string

var (
	NotCreated status = "NOT_CREATED"
	OK         status = "OK"
	Diff       status = "DIFF"
)

func getState() (repos State, err error) {
	var (
		appsJson []byte
		backend  backends.Backend
	)

	backend = backends.GetBackend()

	appsJson, err = backend.GetState()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(appsJson, &repos)
	if err != nil {
		return nil, err
	}

	return repos, nil
}

func checkAppStatus() (appStatuses []AppStatus, err error) {
	var (
		state        State
		templateHash string
		configHash   string
		repoLocation string
		appStatus    AppStatus
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

			appStatus = AppStatus{
				App:              app,
				TemplateLocation: templateLocation,
				ConfigLocation:   configLocation,
				RepoURL:          repo.RepoURL,
			}

			if app.StackID == "" {
				appStatus.Status = NotCreated
				appStatuses = append(appStatuses, appStatus)
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
				appStatus.Status = OK
			} else {
				appStatus.Status = Diff
			}

			appStatuses = append(appStatuses, appStatus)
		}
	}

	return appStatuses, nil
}

func applyStateUpdate(appStatus AppStatus, opOutput aws.CreateUpdateOutput) (err error) {
	var (
		state        State
		newState     State
		newAppState  []App
		templateHash string
		configHash   string
	)

	if state, err = getState(); err != nil {
		return err
	}

	for _, repo := range state {
		if repo.RepoURL == appStatus.RepoURL {
			for _, app := range repo.Apps {
				if app.Name == appStatus.App.Name {
					templateHash, err = utils.GetFileHash(appStatus.TemplateLocation)
					utils.Must(err)

					configHash, err = utils.GetFileHash(appStatus.ConfigLocation)
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

	err = updateState(newState)

	return err
}

func updateState(newState State) (err error) {
	var (
		newStateJSON []byte
		backend      backends.Backend
	)

	backend = backends.GetBackend()

	if err = backend.UpdateState(newStateJSON); err != nil {
		return err
	}

	return err
}

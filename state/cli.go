package state

import (
	"errors"
	"fmt"

	"github.com/craigjames16/cfstate/aws"
	"github.com/craigjames16/cfstate/utils"
	"github.com/urfave/cli/v2"
)

func AddApp(c *cli.Context) (err error) {
	var (
		appName          string
		templateLocation string
		configLocation   string
		repoName         string
		state            []Repo
		newState         []Repo
	)

	utils.ExistOrPrompt("app-name", &appName, c)
	utils.ExistOrPrompt("template-location", &templateLocation, c)
	utils.ExistOrPrompt("config-location", &configLocation, c)
	utils.ExistOrPrompt("repo", &repoName, c)

	if templateLocation[0:1] != "/" {
		templateLocation = "/" + templateLocation
	}

	if configLocation[0:1] != "/" {
		configLocation = "/" + configLocation
	}

	if state, err = getState(); err != nil {
		return err
	}

	for _, repo := range state {
		if repo.RepoName == repoName {
			repo.Apps = append(repo.Apps, App{
				Name:     appName,
				Template: templateLocation,
				Config:   configLocation,
			})
		}

		newState = append(newState, repo)
	}

	if err = updateState(newState); err != nil {
		return err
	}

	return err
}

func SyncState(c *cli.Context) (err error) {
	var (
		opOutput aws.CreateUpdateOutput
	)

	statuses, err := checkAppStatus()
	if err != nil {
		return err
	}

	for _, appStatus := range statuses {
		switch appStatus.Status {
		case NotCreated:
			opOutput, err = aws.CreateStack(aws.AppInput{
				Name:     appStatus.App.Name,
				Template: appStatus.TemplateLocation,
				Config:   appStatus.ConfigLocation,
			})
		case OK:
			fmt.Println("OK")
			continue
		case Diff:
			opOutput, err = aws.UpdateStack(aws.AppInput{
				StackID:  appStatus.App.StackID,
				Template: appStatus.TemplateLocation,
				Config:   appStatus.ConfigLocation,
			})
		}

		if err != nil {
			return err
		}

		utils.Must(applyStateUpdate(appStatus, opOutput))

	}

	return nil
}

func CheckStatus(c *cli.Context) error {
	var (
		err error
	)

	states, err := checkAppStatus()
	if err != nil {
		return err
	}

	for _, state := range states {
		fmt.Printf("%s: %s\n", state.App.Name, state.Status)
	}

	return nil
}

func AddRepo(c *cli.Context) (err error) {
	var (
		repoName string
		repoURL  string
		state    State
		newState State
	)

	utils.ExistOrPrompt("repo-name", &repoName, c)
	utils.ExistOrPrompt("repo-url", &repoURL, c)

	if state, err = getState(); err != nil {
		return err
	}

	for _, repo := range state {
		if repo.RepoName == repoName {
			return errors.New("repo already exists")
		}
	}

	newState = append(state, Repo{
		RepoName: repoName,
		RepoURL:  repoURL,
		Apps:     []App{},
	})

	if err = updateState(newState); err != nil {
		return err
	}

	return err
}

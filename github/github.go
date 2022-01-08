package github

import (
	"fmt"
	"os"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

var REPO_BASE string = "/tmp"

func GetRepo(url string) (location string, err error) {
	var (
		token       string = os.Getenv("GITHUB_ACCESS_TOKEN")
		urlParts    []string
		tmpLocation string
		r           *git.Repository
		ref         *plumbing.Reference
	)

	now := time.Now()
	sec := now.Unix()
	urlParts = strings.Split(url, "/")

	if len(urlParts) < 2 {
		return location, fmt.Errorf("invalid repo url")
	}

	tmpLocation = fmt.Sprintf("%s/%d", REPO_BASE, sec)

	r, err = git.PlainClone(tmpLocation, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: "abc123", // yes, this can be anything except an empty string
			Password: token,
		},
		URL: url,
	})

	if err != nil {
		return tmpLocation, err
	}

	// ... retrieving the branch being pointed by HEAD
	if ref, err = r.Head(); err != nil {
		fmt.Println(err)
		return tmpLocation, err
	}

	location = fmt.Sprintf("%s/%s", REPO_BASE, ref.Hash().String())

	err = os.Rename(tmpLocation, location)
	if err != nil {
		fmt.Println("Repo commit exists locally")
		return location, nil
	}

	return location, err
}

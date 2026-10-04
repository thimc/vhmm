package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	// Name of the locally stored repository
	repoFile = "repo.json"

	// Time after which vhmm considers the repo to be out of date
	// and needs to be redownloaded
	repoDelta = time.Hour
)

type Repository struct {
	URL     string
	Plugins []Plugin
}

func (r *Repository) UpdateNeeded() bool {
	fi, err := os.Stat(repoFile)
	if err != nil {
		return true
	}
	return time.Since(fi.ModTime()) > repoDelta
}

func (r *Repository) Update() error {
	resp, err := http.Get(r.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("repository returned HTTP status %v", resp.Status)
	}
	out, err := os.OpenFile(repoFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	return nil
}

func NewRepository(url string) (Repository, error) {
	repo := Repository{URL: url}
	if repo.UpdateNeeded() {
		fmt.Printf("Downloading new repository from %s to %s\n", url, repoFile)
		if err := repo.Update(); err != nil {
			return repo, err
		}
	}
	f, err := os.Open(repoFile)
	if err != nil {
		return repo, err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&repo.Plugins); err != nil {
		return repo, err
	}
	return repo, nil
}

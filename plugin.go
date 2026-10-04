package main

import (
	"errors"
	"sort"
	"strings"
	"time"
)

type PluginVersion struct {
	VersionNumber string   `json:"version_number"`
	DownloadURL   string   `json:"download_url"`
	Dependencies  []string `json:"dependencies"`
	IsActive      bool     `json:"is_active"`
	UUID4         string   `json:"uuid4"`
}

type Plugin struct {
	Name        string          `json:"name"`
	FullName    string          `json:"full_name"`
	Owner       string          `json:"owner"`
	PackageURL  string          `json:"package_url"`
	DateCreated time.Time       `json:"date_created"`
	DateUpdated time.Time       `json:"date_updated"`
	Versions    []PluginVersion `json:"versions"`
}

func (p *Plugin) Latest() (PluginVersion, error) {
	var av []PluginVersion
	for _, v := range p.Versions {
		if v.IsActive {
			av = append(av, v)
		}
	}
	if len(av) == 0 {
		return PluginVersion{}, errors.New("no active versions available")
	}
	sort.Slice(av, func(i, j int) bool {
		return compareVersions(av[i].VersionNumber, av[j].VersionNumber) > 0
	})
	return av[0], nil
}

func pluginKey(owner, name string) string {
	return strings.ToLower(strings.TrimSpace(name))
	// return strings.ToLower(strings.TrimSpace(owner) + "-" + strings.TrimSpace(owner))
}

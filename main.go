package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
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

type Config struct {
	GameDir       string `json:"game_dir"`
	RepositoryURL string `json:"repository_url"`
	Script        string `json:"script"`
}

func loadConfig() (Config, error) {
	var cfg Config
	data, err := os.ReadFile("config.json")
	if err != nil {
		return cfg, err
	}
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}

func main() {
	var (
		gameDir       = flag.String("game", "", "Path to the root game directory")
		repositoryURL = flag.String("repository", "https://thunderstore.io/c/valheim/api/v1/package/", "URL of the repository")
		dryRun        = flag.Bool("dry-run", false, "Only display changes")
		script        = flag.String("script", "", "Optional; if supplied and vhmm finishes without any issues, run script after")
	)
	flag.Parse()

	cfg, err := loadConfig()
	if err != nil {
		if *gameDir == "" || *repositoryURL == "" {
			flag.Usage()
			os.Exit(1)
		}
		cfg = Config{GameDir: *gameDir, Script: *script}
	}
	if *repositoryURL != "" {
		cfg.RepositoryURL = *repositoryURL
	}
	repo, err := NewRepository(cfg.RepositoryURL)
	if err != nil {
		fmt.Printf("Failed to initialize repository: %v\n", err)
		os.Exit(1)
	}
	var m = NewManager(repo, cfg, *dryRun)
	if err := m.Scan(); err != nil {
		fmt.Printf("Failed to scan installed ps: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d installed plugins\n", len(m.Installed))
	for k, v := range m.Installed {
		fmt.Printf("Found %s\n", k)
		fmt.Printf("- Version:    %s\n", v.Manifest.VersionNumber)
		fmt.Printf("- Directory:  %s\n", v.Directory)
	}
	if err := m.Update(); err != nil {
		fmt.Printf("Failed to update: %v\n", err)
		os.Exit(1)
	}
	if !*dryRun && cfg.Script != "" {
		fmt.Printf("Launching %s\n", cfg.Script)
		cmd := exec.Command("/bin/bash", cfg.Script)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "script failed: %v\n", err)
			os.Exit(1)
		}
	}
}

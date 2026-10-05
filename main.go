package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	gameDir       = flag.String("game", "", "Path to the root game directory")
	repositoryURL = flag.String("repository", "https://thunderstore.io/c/valheim/api/v1/package/", "URL of the repository")
	dryRun        = flag.Bool("dry-run", false, "Only display changes")
)

func main() {
	flag.Parse()
	cfg, err := loadConfig()
	if err != nil {
		if *gameDir == "" || *repositoryURL == "" {
			flag.Usage()
			os.Exit(1)
		}
		cfg = Config{GameDir: *gameDir}
	}
	if *repositoryURL != "" {
		cfg.RepositoryURL = *repositoryURL
	}
	m, err := NewManager(cfg, *dryRun)
	if err != nil {
		fmt.Printf("Failed to initialize manager: %v\n", err)
		os.Exit(1)
	}
	if err := m.Scan(); err != nil {
		fmt.Printf("Failed to scan installed plugin: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d installed plugin(s)\n", len(m.Installed))
	for p, ip := range m.Installed {
		fmt.Printf("%s %s: %s\n", p, ip.Manifest.VersionNumber, ip.Directory)
	}
	if err := m.Update(); err != nil {
		fmt.Printf("Failed to update: %v\n", err)
		os.Exit(1)
	}
}

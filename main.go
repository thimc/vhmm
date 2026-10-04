package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

var (
	gameDir       = flag.String("game", "", "Path to the root game directory")
	repositoryURL = flag.String("repository", "https://thunderstore.io/c/valheim/api/v1/package/", "URL of the repository")
	dryRun        = flag.Bool("dry-run", false, "Only display changes")
	script        = flag.String("script", "", "Optional; if supplied and vhmm finishes without any issues, run script after")
)

func main() {
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
	m, err := NewManager(cfg, *dryRun)
	if err != nil {
		fmt.Printf("Failed to initialize manager: %v\n", err)
		os.Exit(1)
	}
	if err := m.Scan(); err != nil {
		fmt.Printf("Failed to scan installed ps: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d installed plugins\n", len(m.Installed))
	for p, ip := range m.Installed {
		fmt.Printf("%s %s: %s\n", p, ip.Manifest.VersionNumber, ip.Directory)
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

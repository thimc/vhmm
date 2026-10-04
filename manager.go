package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Manager struct {
	GameDir    string
	Repository map[string]Plugin
	Installed  map[string]InstalledPlugin
	Processing map[string]bool
	Processed  map[string]bool
	DryRun     bool
	c          *http.Client
}

func NewManager(r Repository, cfg Config, dryRun bool) *Manager {
	m := &Manager{
		GameDir:    filepath.Join(cfg.GameDir, "BepInEx", "plugins"),
		Repository: make(map[string]Plugin),
		Installed:  make(map[string]InstalledPlugin),
		Processing: make(map[string]bool),
		Processed:  make(map[string]bool),
		DryRun:     dryRun,
		c:          &http.Client{Timeout: 10 * time.Second},
	}
	for _, p := range r.Plugins {
		mk := pluginKey(p.Owner, p.Name)
		m.Repository[mk] = p
	}
	return m
}

type InstalledPlugin struct {
	ManifestPath string
	Directory    string
	Manifest     InstalledManifest
}

type InstalledManifest struct {
	Name          string   `json:"name"`
	VersionNumber string   `json:"version_number"`
	Dependencies  []string `json:"dependencies"`
}

func installedManifestKey(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func dependencyKey(dependency string) string {
	p := strings.Split(strings.TrimSpace(dependency), "-")
	if len(p) < 2 {
		return strings.ToLower(strings.TrimSpace(dependency))
	}
	return strings.ToLower(p[1])
	// return strings.ToLower(strings.Join(parts[:len(p)-1], "-"))
}

func (m *Manager) Scan() error {
	return filepath.Walk(m.GameDir, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() || fi.Name() != "manifest.json" {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var im InstalledManifest
		if err := json.Unmarshal(b, &im); err != nil {
			fmt.Printf("Invalid manifest %s: %v\n", path, err)
			return nil
		}
		if im.Name == "" || im.VersionNumber == "" {
			fmt.Printf("Invalid manifest: \n", path)
			return nil
		}
		k := installedManifestKey(im.Name)
		m.Installed[k] = InstalledPlugin{
			ManifestPath: path,
			Directory:    filepath.Dir(path),
			Manifest:     im,
		}
		return nil
	})
}

func (m *Manager) Update() error {
	var plugins []string
	for m, _ := range m.Installed {
		plugins = append(plugins, m)
	}
	sort.Strings(plugins)
	for _, p := range plugins {
		if err := m.Ensure(p); err != nil {
			fmt.Printf("Failed to ensure p %s: %v\n", p, err)
			return err
		}
	}
	return nil
}

func (m *Manager) Ensure(p string) error {
	if m.Processed[p] {
		fmt.Printf("Already processed: %s\n", p)
		return nil
	}
	if m.Processing[p] {
		fmt.Printf("Already being processed: %s\n", p)
		return nil
	}
	if p == "bepinexpack_valheim" {
		fmt.Printf("Skipping %s\n", p)
		return nil
	}
	rp, exists := m.Repository[p]
	if !exists {
		return fmt.Errorf("%v does not exist in the repository", p)
	}
	lv, err := rp.Latest()
	if err != nil {
		return fmt.Errorf("%s: %w", rp.FullName, err)
	}
	m.Processing[p] = true
	defer delete(m.Processing, p)
	for _, dep := range lv.Dependencies {
		dk := dependencyKey(dep)
		if _, exists := m.Repository[dk]; !exists {
			return fmt.Errorf("dependency %q of %s does not exist in the repository", dep, rp.FullName)
		}
		if err := m.Ensure(dk); err != nil {
			return err
		}
	}
	ik, exists := m.Installed[p]
	if !exists {
		fmt.Printf("Installing %s %s\n", rp.Name, lv.VersionNumber)
		if !m.DryRun {
			if err := m.Install(rp, lv); err != nil {
				return err
			}
		}
	} else {
		switch compareVersions(ik.Manifest.VersionNumber, lv.VersionNumber) {
		case -1:
			fmt.Printf("Updating %s from %s to %s\n", rp.Name, ik.Manifest.VersionNumber, lv.VersionNumber)
			if !m.DryRun {
				if err := m.UpdatePlugin(ik, rp, lv); err != nil {
					return err
				}
			}
		default:
			fmt.Printf("%s is up-to-date (%s)\n", rp.Name, lv.VersionNumber)
		}
	}
	if !m.DryRun {
		// TODO: This consumes unnecessary CPU cycles,
		// Instead of scanning ALL directories we should
		// really just update the current m.Installed[..]
		if err := m.Scan(); err != nil {
			return err
		}
	}
	m.Processed[p] = true
	return nil
}

func (m *Manager) download(url string, out io.Writer) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := m.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %s", resp.Status)
	}
	_, err = io.Copy(out, resp.Body)
	return err
}

func (m *Manager) extract(path, dest string) error {
	fmt.Printf("Extracting %s to %s\n", path, dest)
	r, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		tp := filepath.Join(dest, f.Name)
		fmt.Printf("-> %s\n", tp)
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(tp, 0755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(tp), 0755); err != nil {
			return err
		}
		in, err := f.Open()
		if err != nil {
			return err
		}
		out, err := os.OpenFile(tp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode().Perm())
		if err != nil {
			in.Close()
			return err
		}
		if _, err = io.Copy(out, in); err != nil {
			return err
		}
		if err = in.Close(); err != nil {
			return err
		}
		if err = out.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) Install(p Plugin, pv PluginVersion) error {
	if pv.DownloadURL == "" {
		return fmt.Errorf("%s is missing a download URL", p.FullName)
	}
	tf, err := os.CreateTemp("", "plugin-*.zip")
	if err != nil {
		return err
	}
	tp := tf.Name()
	target := fmt.Sprintf("%s-%s", p.FullName, pv.VersionNumber)
	defer os.Remove(tp)
	if err := m.download(pv.DownloadURL, tf); err != nil {
		tf.Close()
		return fmt.Errorf("downloading %s: %w", p.FullName, err)
	}
	if err := tf.Close(); err != nil {
		return err
	}
	fmt.Printf("downloaded to %s\n", tp)
	if err := m.extract(tp, filepath.Join(m.GameDir, target)); err != nil {
		return fmt.Errorf("extracting %s: %w", p.FullName, err)
	}
	return nil
}

func (m *Manager) UpdatePlugin(i InstalledPlugin, p Plugin, v PluginVersion) error {
	td := i.Directory
	fmt.Printf("Removing old installation of %s: %s\n", p.FullName, td)
	if err := os.RemoveAll(td); err != nil {
		return err
	}
	return m.Install(p, v)
}

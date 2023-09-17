//go:build mage
// +build mage

// Package mage contains setup for magefiles.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Aliases for some of the more used commands.
var Aliases = map[string]interface{}{
	"t:u":   Test.Unit,
	"t:ua":  Test.UnitAll,
	"t:uw":  Test.UnitWatch,
	"t:uwa": Test.UnitWatchAll,
	"l":     Lint,
}

// Dependencies downloads/upgrades dependencies.
func Dependencies() error {
	deps := []string{
		"github.com/cespare/reflex@latest",
		"github.com/golangci/golangci-lint/cmd/golangci-lint@latest",
		"github.com/psampaz/go-mod-outdated@latest",
		"github.com/jondot/goweight@latest",
		"golang.org/x/vuln/cmd/govulncheck@latest",
	}
	fmt.Println("Installing Deps...")
	for _, d := range deps {
		err := sh.Run("go", "install", d)
		if err != nil {
			return fmt.Errorf("installing %s: %w", d, err)
		}
	}
	deps = []string{
		"golang.org/x/tools/cmd/cover",
		"github.com/sonatype-nexus-community/nancy@latest",
	}
	for _, d := range deps {
		err := sh.Run("go", "get", "-t", "-u", d)
		if err != nil {
			return fmt.Errorf("installing %s: %w", d, err)
		}
	}
	err := sh.Run("go", "get", "-d", "-u", "./...")
	if err != nil {
		return fmt.Errorf("updating all dependencies: %w", err)
	}
	return Tidy()
}

// Run runs the game.
func Run() error {
	return sh.RunV("go", "run", ".")
}

// Lint lints the code.
func Lint() error {
	err := sh.RunV("go", "fmt", "./...")
	if err != nil {
		return fmt.Errorf("running go fmt tool: %w", err)
	}
	err = sh.RunV("go", "vet", "./...")
	if err != nil {
		return fmt.Errorf("running go vet tool: %w", err)
	}
	err = sh.RunV("golangci-lint", "run", "./...")
	if err != nil {
		return fmt.Errorf("running golangci-lint: %w", err)
	}
	return nil
}

// Clean tidies the mod file and cleans the test cache.
func Clean() error {
	mg.Deps(Tidy)
	return sh.Run("go", "clean", "-testcache")
}

// Tidy tidies the modules.
func Tidy() error {
	return sh.RunV("go", "mod", "tidy")
}

// Test is a namespace for running tests.
type Test mg.Namespace

// Unit runs the unit tests.
func (Test) Unit() error {
	_ = Tidy()
	return sh.RunV("go", "test", "-trimpath", "-failfast", "-short", "./...")
}

// UnitAll runs the unit tests without running without skipping slow tests.
func (Test) UnitAll() error {
	_ = Tidy()
	return sh.RunV("go", "test", "-trimpath", "-failfast", "./...")
}

// UnitWatch watches for file changes and runs the unit tests.
func (t Test) UnitWatch(ctx context.Context) error {
	ch, err := watchChanges(ctx)
	if err != nil {
		return err
	}
	_ = t.Unit()
	fmt.Println(strings.Repeat("#", 40))
	for range ch {
		err := t.Unit()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(strings.Repeat("#", 40))
	}
	return nil
}

// UnitWatchAll watches for file changes and runs the unit tests without
// skipping slow tests.
func (t Test) UnitWatchAll(ctx context.Context) error {
	ch, err := watchChanges(ctx)
	if err != nil {
		return err
	}
	_ = t.UnitAll()
	fmt.Println(strings.Repeat("#", 40))
	for range ch {
		err := t.UnitAll()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(strings.Repeat("#", 40))
	}
	return nil
}

// CI runs all tests, used for github actions.
func (Test) CI() error {
	mg.Deps(Tidy)
	return sh.RunV("go", "test", "-trimpath", "-failfast", "-race", "./...")
}

func watchChanges(ctx context.Context) (chan struct{}, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	ch := make(chan struct{})
	go func() {
		<-ctx.Done()
		close(ch)
		_ = watcher.Close()
	}()

	err = watchFiles(watcher)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			// watching for new directories.
			_ = watchFiles(watcher)
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				ext := filepath.Ext(event.Name)
				if !slices.Contains(exts, ext) {
					continue
				}
				if event.Op == fsnotify.Write {
					select {
					case ch <- struct{}{}:
					default:
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	return ch, nil
}

var exts = []string{
	".go",
	".mod",
}

func watchFiles(watcher *fsnotify.Watcher) error {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walking directory: %w", err)
		}
		if !info.IsDir() {
			return nil
		}
		err = watcher.Add(path)
		if err != nil {
			log.Println("error:", err)
		}
		return nil
	})
	return err
}

// Audit audits the code for updates, vulnerabilities and binary weight.
func Audit() error {
	err := sh.RunV("govulncheck", "./...")
	if err != nil {
		return fmt.Errorf("running govulncheck: %w", err)
	}

	err = pipe(
		[]string{"go", "list", "-u", "-m", "-json", "all"},
		[]string{"go-mod-outdated", "-update", "-direct"},
	)
	if err != nil {
		return fmt.Errorf("getting update list: %w", err)
	}

	err = pipe(
		[]string{"go", "list", "-json", "-deps"},
		[]string{"nancy", "sleuth"},
	)
	if err != nil {
		return fmt.Errorf("getting nancy slueth: %w", err)
	}

	out, err := sh.Output("goweight")
	if err != nil {
		return fmt.Errorf("getting package list: %w", err)
	}
	split := strings.Split(out, "\n")
	length := min(len(split), 20)
	fmt.Println(strings.Join(split[:length], "\n"))

	return nil
}

func pipe(cmd1, cmd2 []string) error {
	fmt.Printf("Running %s | %s\n", strings.Join(cmd1, " "), strings.Join(cmd2, " "))
	out, err := sh.Output(cmd1[0], cmd1[1:]...)
	if err != nil {
		fmt.Println(out)
		return fmt.Errorf("getting package list: %w", err)
	}

	c := exec.Command(cmd2[0], cmd2[1:]...)
	c.Env = os.Environ()
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	c.Stdin = strings.NewReader(out)

	return c.Run()
}

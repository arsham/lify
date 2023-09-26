//go:build mage
// +build mage

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/magefile/mage/mg"
)

var (
	osxPairs = [][3]string{
		{"darwin", "amd64"},
		{"darwin", "arm"},
		{"darwin", "arm64"},
	}
	linuxPairs = [][3]string{
		{"linux", "amd64"},
		{"linux", "arm"},
		{"linux", "arm64"},
	}
	winPairs = [][3]string{
		{"windows", "386", ".exe"},
		{"windows", "amd64", ".exe"},
	}
	jsPairs = [][3]string{
		{"js", "wasm"},
	}
	androidPairs = [][3]string{
		{"android", "arm"},
	}
	archPairFlags = map[[3]string][]string{
		{"windows", "386", ".exe"}:   {"-ldflags=-H=windowsgui"},
		{"windows", "amd64", ".exe"}: {"-ldflags=-H=windowsgui"},
	}
	cgoEnv = map[string][]string{
		"darwin": {
			// "CC=arm-linux-gnueabihf-gcc",
			"CGO_ENABLED=1",
		},
	}
)

// Build namespace for general build tooling.
type Build mg.Namespace

// All distros should be built.
func (Build) All() {
	mg.SerialDeps(
		Tidy,
		Build.Linux,
		// Build.OSX,
		Build.JS,
		Build.Windows,
		// Build.Android,
	)
}

// Linux builds for linux.
func (Build) Linux() error {
	return buildWithPairs(linuxPairs)
}

// OSX builds for OSX.
func (Build) OSX() error {
	return buildWithPairs(osxPairs)
}

// Windows builds for windows.
func (Build) Windows() error {
	return buildWithPairs(winPairs)
}

// JS builds for javascript (wasm).
func (Build) JS() error {
	return buildWithPairs(jsPairs)
}

// Android builds for android.
func (Build) Android() error {
	return buildWithPairs(androidPairs)
}

func buildWithPairs(pairs [][3]string) error {
	if err := os.Chdir("build"); err != nil {
		return fmt.Errorf("changing dir: %w", err)
	}
	defer func() {
		if err := os.Chdir("../"); err != nil {
			panic(fmt.Errorf("changing dir: %w", err))
		}
	}()

	for _, pair := range pairs {
		err := os.Setenv("GOOS", pair[0])
		if err != nil {
			return fmt.Errorf("setting GOOS env: %w", err)
		}
		err = os.Setenv("GOARCH", pair[1])
		if err != nil {
			return fmt.Errorf("setting GOARCH env: %w", err)
		}
		if envs, ok := cgoEnv[pair[0]]; ok {
			fmt.Println("CGO is enabled for", pair[0])
			for _, env := range envs {
				e := strings.Split(env, "=")
				err = os.Setenv(e[0], e[1])
				if err != nil {
					return fmt.Errorf("setting CGO env: %w", err)
				}
			}
		}
		buildName := fmt.Sprintf("neuragene_%s_%s%s", pair[0], pair[1], pair[2])
		toRun := []string{"build", "--tags", "prod", "-o", buildName}
		if flags, ok := archPairFlags[pair]; ok {
			toRun = append(toRun, flags...)
		}
		toRun = append(toRun, "github.com/arsham/neuragene")
		fmt.Println("Running: go ", toRun)
		cmd := exec.Command("go", toRun...)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Println("Running command:", err)
		}
		if out.Len() != 0 {
			fmt.Printf("%s\n", out.String())
		}
		fmt.Printf("Build for %s is finished!\n\n", pair[0])
	}
	return nil
}

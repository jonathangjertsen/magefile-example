// +build mage

package main

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	EXE = "example"
)

// Builds an executable for Windows AMD64
func BuildWindowsAmd64() error {
	return build("windows", "amd64", EXE, []string{"-tags", "release"})
}

// Builds an executable for Mac AMD64
func BuildMacAmd64() error {
	return build("darwin", "amd64", EXE, []string{"-tags", "release"})
}

// Builds an executable for Mac ARM64
func BuildMacArm64() error {
	return build("darwin", "arm64", EXE, []string{"-tags", "release"})
}

// Builds an executable for Linux AMD64
func BuildLinuxAmd64() error {
	return build("linux", "amd64", EXE, []string{"-tags", "release"})
}

// Builds an executable for Linux ARM64
func BuildLinuxArm64() error {
	return build("linux", "arm64", EXE, []string{"-tags", "release"})
}

// Builds an executable for this computer
func Build() error {
	return build(runtime.GOOS, runtime.GOARCH, EXE, []string{"-tags", "release"})
}

// Builds a debug executable for this computer
func Debug() error {
	return build(runtime.GOOS, runtime.GOARCH, fmt.Sprintf("%s.debug", EXE), []string{"-tags", "debug"})
}

// Runs the program in debug mode without arguments
func Run() error {
	mg.Deps(Debug)
	output, err := run(
		executablePath(runtime.GOOS, runtime.GOARCH, fmt.Sprintf("%s.debug", EXE)),
		[]string{},
		map[string]string{},
	)
	fmt.Print(output)
	if err != nil {
		return err
	}
	return nil
}

// Runs the program in release mode without arguments
func RunRelease() error {
	mg.Deps(Build)
	output, err := run(executablePath(runtime.GOOS, runtime.GOARCH, EXE), []string{}, map[string]string{})
	fmt.Print(output)
	if err != nil {
		return err
	}
	return nil
}

// Builds an executable for all supported platforms
func BuildAll() {
	parallelBuild([](func() error){
		BuildWindowsAmd64,
		BuildMacAmd64,
		BuildMacArm64,
		BuildLinuxAmd64,
		BuildLinuxArm64,
		Debug,
	})
}

// Runs everything that a CI system might want to do
func Ci() {
	mg.Deps(Check)
	mg.Deps(CheckRepoClean)
	mg.Deps(Test)
	mg.Deps(BuildAll)
	mg.Deps(Run)
	mg.Deps(RunRelease)
}

// Runs go vet and go fmt, and checks that they don't say anything
func Check() error {
	output, err := run("go", []string{"vet", "./..."}, map[string]string{})
	if err != nil {
		return err
	}
	if output != "" {
		return fmt.Errorf("go vet says something:\n%s", output)
	}

	output, err = run("go", []string{"fmt", "./..."}, map[string]string{})
	if err != nil {
		return err
	}
	if output != "" {
		return fmt.Errorf("go fmt says something:\n%s", output)
	}
	return nil
}

// Checks that the repo is clean
func CheckRepoClean() error {
	output, err := run("git", []string{"status", "--porcelain"}, map[string]string{})
	if err != nil {
		return err
	}
	if output != "" {
		return fmt.Errorf("git status --porcelain says something:\n%s", output)
	}
	return nil
}

// Runs go test in verbose mode and prettifies the output
func Test() error {
	output, err := run("go", []string{"test", "-v", "./..."}, map[string]string{})
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "[no test files]") {
			continue
		}
		if strings.HasPrefix(line, "ok") {
			color.HiGreen(line)
		} else if strings.Contains(line, "FAIL") {
			color.HiRed(line)
		} else {
			color.New().Println(line)
		}
	}
	return err
}

// Cleans the bin directory
func Clean() error {
	fmt.Println("Removing bin")
	return sh.Rm("bin")
}

func parallelBuild(builders [](func() error)) {
	var wg sync.WaitGroup

	for _, builder := range builders {
		wg.Add(1)
		go (func(builder func() error, wg *sync.WaitGroup) {
			defer wg.Done()
			builder()
		})(builder, &wg)
	}
	wg.Wait()
}

func executablePath(os, arch, executable string) string {
	extension := ""
	if os == "windows" {
		extension = ".exe"
	}
	return fmt.Sprintf("bin/%s-%s/%s%s", os, arch, executable, extension)
}

func build(os, arch, executable string, args []string) error {
	command := []string{
		"build",
		"-o", executablePath(os, arch, executable),
	}
	command = append(command, args...)
	command = append(command, ".")
	output, err := run("go", command, map[string]string{
		"GOOS":   os,
		"GOARCH": arch,
	})
	fmt.Print(output)
	return err
}

func run(program string, args []string, env map[string]string) (string, error) {
	// Make string representation of command
	fullArgs := append([]string{program}, args...)
	cmdStr := strings.Join(fullArgs, " ")

	// Make string representation of environment
	envStrBuf := new(bytes.Buffer)
	for key, value := range env {
		fmt.Fprintf(envStrBuf, "%s=\"%s\", ", key, value)
	}
	envStr := string(bytes.TrimRight(envStrBuf.Bytes(), ", "))

	// Show info
	fmt.Println("Running '" + cmdStr + "'" + " with env " + envStr)

	// Run
	return sh.OutputWith(env, program, args...)
}

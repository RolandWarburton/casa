package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

// returns the exit code or -1 if the exit code cannot be determined
// intended for use with spawned shell commands
func checkExitErrorCode(err error) int {
	if exitError, ok := err.(*exec.ExitError); ok {
		exitCode := exitError.ExitCode()
		fmt.Printf("Command exited with code: %d\n", exitCode)
		return exitCode
	} else {
		log.Fatal(err)
		return -1
	}
}

// pass in (command, true) to enable debugging
func spawnShellCommand(command string, params ...bool) int {
	isDebug := len(params) > 0 && params[0]

	shell := []string{"/usr/bin/sh", "-c"}
	executable := append(shell, command)
	cmd := exec.Command(executable[0], executable[1:]...)

	// Create a pipe for the command's stdout
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	// Start the command
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	// Read the stdout asynchronously
	output, err := ioutil.ReadAll(stdoutPipe)
	if err != nil {
		exitCode := checkExitErrorCode(err)
		if exitCode != -1 {
			return exitCode
		} else {
			log.Fatal(err)
			return 1
		}
	}

	// Wait for the command to finish
	err = cmd.Wait()
	if err != nil {
		exitCode := checkExitErrorCode(err)
		if exitCode != -1 {
			return exitCode
		} else {
			log.Fatal(err)
			return 1
		}
	}

	if isDebug {
		fmt.Printf("Command output:\n%s\n", string(output))
	}
	return 0
}

func createSymlink(source, target string) error {
	// Expand the tilde in the target path to the user's home directory
	if strings.HasPrefix(target, "~") {
		usr, err := user.Current()
		if err != nil {
			return fmt.Errorf("[ERROR] failed to get current user: %v", err)
		}
		target = filepath.Join(usr.HomeDir, target[1:])
	}

	// Get the parent directory of the target
	parentDir := filepath.Dir(target)

	// Check if the parent directory exists
	_, err := os.Stat(parentDir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("parent directory does not exist: %s", parentDir)
		}
		return fmt.Errorf("[ERROR] failed to check parent directory: %v", err)
	}

	// the source is relative to the working directory
	// resolve the working directory to get the full source path
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get working directory")
	}
	source = path.Join(wd, source)

	// Create the symlink
	if err := os.Symlink(source, target); err != nil {
		return fmt.Errorf("[ERROR] failed to create symlink: %v", err)
	}

	fmt.Printf("Created symlink: %s -> %s\n", target, source)
	return nil
}

func LinkFile(link Link) error {
	// check if we should run this link
	if len(link.If) > 0 {
		if exitCode := spawnShellCommand(link.If); exitCode != 0 {
			fmt.Printf("skipping link %s -> %s\n", link.Path, link.Destination)
			return nil
		}
	}

	// create the symlink
	err := createSymlink(link.Path, link.Destination)
	if err != nil {
		return fmt.Errorf("[ERROR] failed to create symlink: %v", err)
	}
	return nil
}

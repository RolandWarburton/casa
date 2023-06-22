package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SymLinkResult struct {
	Source, Target string
	Err            error
}

func walkSymLinks(path string, info os.FileInfo, err error, ch chan<- SymLinkResult) error {
	if err != nil {
		ch <- SymLinkResult{Source: path, Target: "", Err: fmt.Errorf("error walking symlinks %v", err)}
	}

	// Check if the file is a symlink
	if info.Mode()&os.ModeSymlink == os.ModeSymlink {
		// Get the absolute path of the symlink target
		target, err := os.Readlink(path)
		if err != nil {
			ch <- SymLinkResult{Source: path, Target: "", Err: fmt.Errorf("failed to read symlink target: %v", err)}
		}

		resolvedTarget, err := filepath.EvalSymlinks(target)
		if err != nil {
			// we should take action because the symlink target could not be evaluated
			// so the source should be deleted (the symlink file)
			ch <- SymLinkResult{Source: path, Target: target, Err: nil}
		} else {
			// since the resolved target was found we can pass back an error to skip it
			ch <- SymLinkResult{Source: path, Target: resolvedTarget, Err: fmt.Errorf("got symlink target %s -> %s (skipping)", path, resolvedTarget)}
		}
	}
	return nil
}

func CleanSymlinks(path string) ([]string, error) {
	ch := make(chan SymLinkResult)

	go func() {
		defer close(ch)
		err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			return walkSymLinks(p, info, err, ch)
		})

		if err != nil {
			fmt.Printf("error walking path: %v\n", err)
		}
	}()

	// Collect the symlinks from the channel
	var symlinks []string
	for symlink := range ch {
		// skip symlinks that failed for whatever reason
		if symlink.Err != nil {
			continue
		}
		symlinks = append(symlinks, symlink.Target)
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("could not get working directory")
		}
		// if the target of the symlink is pointing to a file in the working directory
		// which is where dotfiles will be housed. Then we can remove that link
		fmt.Printf("%s -> %s\n", symlink.Source, symlink.Target)
		if strings.HasPrefix(symlink.Target, wd) {
			err := os.Remove(symlink.Source)
			if err != nil {
				return nil, fmt.Errorf("failed to remove symlink: %v", err)
			} else {
				fmt.Printf("[CLEAN] cleaned %s\n", path)
			}
		}
	}

	return symlinks, nil
}

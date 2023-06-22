package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Link struct {
	Destination string `yaml:"destination"`
	Path        string `yaml:"path"`
	If          string `yaml:"if"`
}

type Dotfiles struct {
	Link   []Link   `yaml:"link"`
	Create []string `yaml:"create"`
	Clean  []string `yaml:"clean"`
}

func Program(yamlFile string) error {
	yamlData, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %v", err)
	}

	var config Dotfiles
	err = yaml.Unmarshal(yamlData, &config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %v", err)
	}

	// clean all the required folders
	for _, cleanPath := range config.Clean {
		_, err := CleanSymlinks(cleanPath)
		if err != nil {
			fmt.Println(err)
		}
	}

	// create all the required folders
	for _, createPath := range config.Create {
		err := CreatePathDir(createPath, true)
		if err != nil {
			fmt.Println(err)
		}
	}

	// create all the required links
	for _, link := range config.Link {
		err := LinkFile(link)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

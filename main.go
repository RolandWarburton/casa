package main

import (
	"io/ioutil"
	"log"

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
}

func main() {
	yamlFile := "install.conf.yaml"

	yamlData, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}

	var config Dotfiles
	err = yaml.Unmarshal(yamlData, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML: %v", err)
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
}

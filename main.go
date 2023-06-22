package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	var configPath string

	app := &cli.App{
		Name:     "casa",
		HelpName: "casa",
		Usage:    "dotfiles linker",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Value:       "install.conf.yaml",
				Aliases:     []string{"c"},
				Usage:       "configuration file",
				Destination: &configPath,
				Required:    true,
			},
		},
		Action: func(_ *cli.Context) error {
			Program(configPath)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

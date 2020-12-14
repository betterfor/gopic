package cmd

import (
	"github.com/urfave/cli/v2"
	"os"
)

var Config = &cli.Command{
	Name:  "config",
	Usage: "set or get config",
	Description: `
Save uploading config of image gallery, it can create from yaml config".
It save config in your $HOME/.gopic/gopic.yaml, 
if it is not exist, will create file.`,
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "list", Usage: "list all config", Value: false, Destination: &configList},
		&cli.StringFlag{Name: "file", Aliases: []string{"f"}, Usage: "", Value: ""},
	},
	Before: configBefore,
	Action: configRun,
}

var (
	home, _ = os.UserHomeDir()

	configList bool
)

const ()

func configBefore(c *cli.Context) error {

}

func configRun(c *cli.Context) error {

}

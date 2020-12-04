package cmd

import "github.com/urfave/cli/v2"

var Config = &cli.Command{
	Name:  "config",
	Usage: "set or get config",
	Description: `
Save uploading config of image gallery,it use format as "yaml/json/toml/ini".
It save config in your $HOME/.gopic/gopic.yaml, 
if it is not exist, will create file.`,
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "list", Usage: "list ", Value: "example.yaml"},
		&cli.StringFlag{
			Name:  "",
			Usage: "",
			Value: "",
		},
	},
}

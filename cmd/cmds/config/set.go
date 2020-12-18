/**
 *Created by XieJian on 2020/12/18 10:44
 *@Desc:
 */
package config

import (
	"github.com/urfave/cli/v2"
)

var set = &cli.Command{
	Name:   "set",
	Usage:  "set your pic bed setting",
	Action: setConfig,
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "autoRename", Aliases: []string{"r"}, Usage: "set config autoRename to rename uploading file in timestamp", Value: ""},
		&cli.StringFlag{Name: "value", Aliases: []string{"v"}, Usage: "set config details", Value: ""},
		&cli.StringSliceFlag{Name: "file", Aliases: []string{"f"}, Usage: "use file to set config", Value: nil},
	},
}

func setConfig(c *cli.Context) error {
	if auto := c.String("autoRename"); len(auto) != 0 {
		switch auto {
		case "true":
			CoreConfig.Base.AutoRename = true
		case "false":
			CoreConfig.Base.AutoRename = false
		}
	}
	val := c.String("value")
	if len(val) != 0 {

	}
	return nil
}

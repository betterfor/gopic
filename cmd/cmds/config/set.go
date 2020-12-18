/**
 *Created by XieJian on 2020/12/18 10:44
 *@Desc:
 */
package config

import (
	"fmt"
	"github.com/betterfor/gopic/core"
	"github.com/betterfor/gopic/core/plugins"
	"github.com/urfave/cli/v2"
	"reflect"
	"strings"
)

var set = &cli.Command{
	Name:   "set",
	Usage:  "set your pic bed setting",
	Action: setConfig,
	After:  SaveConfig,
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
	// value is like: github.token=xxx
	if flags := strings.Split(c.String("value"), "="); len(flags) == 2 {
		return valueParse(flags[0], flags[1])
	}
	return nil
}

func valueParse(param, value string) error {
	keys := strings.Split(param, ".")
	if len(keys) != 2 {
		return fmt.Errorf("incorrect format: %s", param)
	}
	t, v := keys[0], keys[1]
	switch t {
	case core.Github:
		var opts = &plugins.GithubOpts{}
		settings := CoreConfig.Settings[core.Github]
		if settings != nil {
			opts = opts.Unmarshal([]byte(settings.String()))
		}
		ref := reflect.TypeOf(*opts)
		refv := reflect.ValueOf(opts).Elem()
		for i := 0; i < ref.NumField(); i++ {
			if ref.Field(i).Tag.Get("json") == v {
				refv.Field(i).SetString(value)
			}
		}
		CoreConfig.Settings[core.Github] = opts
	}
	return nil
}

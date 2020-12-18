/**
 *Created by XieJian on 2020/12/17 16:34
 *@Desc:
 */
package config

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

var list = &cli.Command{
	Name:  "list",
	Usage: "list config use yaml/json format",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "output", Aliases: []string{"o"}, Usage: "Output format. One of: json|yaml|wide|name", Value: "yaml"},
		&cli.StringFlag{Name: "simple", Aliases: []string{"s"}, Usage: "", Value: "yaml"},
	},
	Action: listConfig,
}

func listConfig(c *cli.Context) error {
	if output := c.String("output"); len(output) != 0 {
		switch output {
		case "json":
			bts, _ := json.MarshalIndent(CoreConfig, " ", "  ")
			fmt.Println(string(bts))
		case "yaml", "yml":
			bts, _ := yaml.Marshal(CoreConfig)
			fmt.Println(string(bts))
		case "wide":
			fmt.Printf("%s\t%s", "name", "setting")
			for key, upload := range CoreConfig.Settings {
				fmt.Printf("%s\t%s", key, upload.String())
			}
		case "name":
			for name := range CoreConfig.Settings {
				fmt.Println(name)
			}
		default:
			return fmt.Errorf("unsupport format: %s", output)
		}
	}
	return nil
}

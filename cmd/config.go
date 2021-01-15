package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
)

type configCmd struct {
	out     io.Writer
	listAll bool
	content map[string]string
	use     string
}

func newConfigCmd(out io.Writer) *cobra.Command {
	config := &configCmd{out: out}
	cmd := &cobra.Command{
		Use:   "config",
		Short: "list or set config",
		Run: func(cmd *cobra.Command, args []string) {
			config.run()
		},
	}
	f := cmd.Flags()
	f.BoolVarP(&config.listAll, "list", "a", false, "show all config")
	f.StringToStringVar(&config.content, "set", nil, "save config")
	f.StringVar(&config.use, "use", "", "choose config")
	return cmd
}

func (c *configCmd) run() {
	if c.listAll {
		fmt.Fprintln(c.out, "list config\n", cfg.String())
		return
	}

	if len(c.use) != 0 {
		//cfg.Current = c.use
		viper.Set("current", c.use)
		fmt.Fprintln(c.out, "convert upload target: ", c.use)
		return
	}
	if len(c.content) != 0 {
		for key, val := range c.content {
			viper.Set(key, val)
		}
	}
}

package cmd

import (
	"fmt"
	"github.com/betterfor/gopic/core/plugins/uploader"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"strings"
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
		Short: "configuration",
		Run: func(cmd *cobra.Command, args []string) {
			config.run()
		},
	}

	cmd.AddCommand(
		listConfigCmd(out),
		useConfigCmd(out),
		setConfigCmd(out),
	)
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
			fmt.Println(key, val)
			viper.Set(key, val)
		}
	}
}

func listConfigCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "config list",
		Run: func(cmd *cobra.Command, args []string) {
			var config = cfg
			config.Uploaded = nil
			fmt.Fprintln(out, config.String())
		},
	}
	return cmd
}

func useConfigCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use",
		Short: "config change upload kind",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("cannot known what kind of uploader use")
			}
			curr := uploader.Convert(args[0])
			if curr == uploader.Unknown {
				return fmt.Errorf("unknown kind of %s", args[0])
			}
			viper.Set("current", curr)
			fmt.Fprintln(out, "current upload target: ", curr)
			return nil
		},
	}
	return cmd
}

func setConfigCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "config set",
		Run: func(cmd *cobra.Command, args []string) {
			for _, arg := range args {
				keys := strings.Split(arg, "=")
				if len(keys) == 2 {
					fmt.Fprintln(out, keys[0], keys[1])
					viper.Set(keys[0], keys[1])
				}
			}
		},
	}
	return cmd
}

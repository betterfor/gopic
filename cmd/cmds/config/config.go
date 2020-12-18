package config

import (
	"github.com/betterfor/gopic/core"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

var Config = &cli.Command{
	Name:  "config",
	Usage: "set or get config",
	Description: `
Save uploading config of image gallery, it can create from yaml config".
It save config in your $HOME/.gopic/gopic.yaml, 
if it is not exist, will create file.`,
	Subcommands: []*cli.Command{list, set},
	After:       SaveConfig,
}

var (
	home, _    = os.UserHomeDir()
	configPath = filepath.Join(home, ".gopic")
	configFile = filepath.Join(home, ".gopic/config.yaml")
	CoreConfig *core.Config
)

func BeforeConfig(c *cli.Context) error {
	err := os.MkdirAll(configPath, os.ModeDir)
	if err != nil {
		return err
	}
	if _, err = os.Stat(configFile); err != nil {
		file, err := os.Create(configFile)
		if err != nil {
			return err
		}
		return file.Close()
	}
	return ParseConfig(c)
}

// parse config file
func ParseConfig(c *cli.Context) error {
	var cfg core.Config
	file, err := os.OpenFile(configFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	err = yaml.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return err
	}
	CoreConfig = &cfg
	return nil
}

// save config
func SaveConfig(c *cli.Context) error {
	bts, err := yaml.Marshal(CoreConfig)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configFile, bts, os.ModePerm)
}

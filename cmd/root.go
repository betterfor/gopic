package cmd

import (
	"fmt"
	"github.com/betterfor/gopic/core"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	cfgFile string
	cfg     = &core.Config{}
	debug   bool
)

var rootCmd = &cobra.Command{
	Use:   "gopic",
	Short: "Manage picture to remote storage",
	Long: `Gopic is a tool for uploading images.
It's easily, quickly, conveniently and support github| ... now.
After your uploading images, you can get a link to save in your blog|markdown|article...'`,
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		return viper.WriteConfigAs(cfgFile)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gopic/config.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Help message for debug")
	out := os.Stdout
	rootCmd.AddCommand(
		newUploadCmd(out),
		newConfigCmd(out))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println("home error:", err)
			os.Exit(1)
		}

		// Search config in home directory with name "config" (without extension).
		viper.AddConfigPath(home + "/.betterfor/gopic")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		cfgPath := filepath.Join(home, ".betterfor", "gopic")
		os.MkdirAll(cfgPath, os.ModeDir)
		cfgFile = filepath.Join(cfgPath, "config.yaml")
	}

	if _, err := os.Stat(cfgFile); err != nil {
		// If cfgFile not exist, create file and write default config.
		ioutil.WriteFile(cfgFile, []byte(cfg.String()), os.ModePerm)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if debug {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}
	viper.Unmarshal(cfg)
}

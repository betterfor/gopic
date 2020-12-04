package main

import (
	"github.com/betterfor/gopic/cmd"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Usage = "gopic is a terminal tool for quickly uploading images and getting URL links to images"
	app.Version = "0.0.1"
	app.Commands = []*cli.Command{
		cmd.Config,
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jwmwalrus/bumpy/internal/task"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "bumpy-ride",
		Version:  "v0.50.1",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "John M",
				Email: "jwmwalrus@gmail.com",
			},
		},
		Copyright: "(c) 2021 WalrusInc Solutions",
		HelpName:  "bumpy",
		Usage:     "A versioning tool",
		UsageText: "bumpy [command] [options ...]",
		ExitErrHandler: func(c *cli.Context, err error) {
			if err != nil {
				fmt.Fprintf(c.App.ErrWriter, err.Error()+"\n")
			}
		},
		Commands: []*cli.Command{
			task.Init(),
			task.Bump(),
			task.Sync(),
			task.Tag(),
			task.Version(),
		},
	}

	app.Run(os.Args)
}

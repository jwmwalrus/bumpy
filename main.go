package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jwmwalrus/bumpy-ride/internal/task"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "bumpy-ride",
		Version:  "v0.1.0",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "John M",
				Email: "jwmwalrus@gmail.com",
			},
		},
		Copyright: "(c) 2021 WalrusInc Solutions",
		HelpName:  "contrive",
		Usage:     "demonstrate available API",
		UsageText: "contrive - demonstrating the available API",
		ArgsUsage: "[args and such]",
		ExitErrHandler: func(c *cli.Context, err error) {
			if err != nil {
				fmt.Fprintf(c.App.ErrWriter, err.Error()+"\n")
			}
		},
		Commands: getTasks(),
	}

	app.Run(os.Args)
}

func getTasks() []*cli.Command {
	return []*cli.Command{
		task.Init(),
		task.Bump(),
		task.Sync(),
		task.Tag(),
	}
}

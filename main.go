package main

import (
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jwmwalrus/bumpy/task"
	"github.com/jwmwalrus/bumpy/version"
	"github.com/urfave/cli/v2"
)

//go:embed version.json
var versionJSON []byte

var appVersion version.Version

var logger *slog.Logger

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	slog.SetDefault(logger)

	app := &cli.App{
		Name:      "bumpy-ride",
		Version:   appVersion.String(),
		Compiled:  time.Now(),
		Copyright: "(c) 2022 WalrusAhead Solutions",
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
			task.Config(),
		},
	}

	app.Run(os.Args)
}

func init() {
	if err := appVersion.Read(versionJSON); err != nil {
		panic(err)
	}
}

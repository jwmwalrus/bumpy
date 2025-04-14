package main

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"os"

	"github.com/jwmwalrus/bumpy/task"
	"github.com/jwmwalrus/bumpy/version"
	"github.com/urfave/cli/v3"
)

//go:embed version.json
var versionJSON []byte

var appVersion version.Version

var logger *slog.Logger

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	slog.SetDefault(logger)

	app := &cli.Command{
		Name:      "bumpy-ride",
		Version:   appVersion.String(),
		Copyright: "(c) 2022 WalrusAhead Solutions",
		Usage:     "A versioning tool",
		UsageText: "bumpy [command] [options ...]",
		ExitErrHandler: func(ctx context.Context, c *cli.Command, err error) {
			if err != nil {
				fmt.Fprintf(c.ErrWriter, err.Error()+"\n")
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

	app.Run(context.Background(), os.Args)
}

func init() {
	if err := appVersion.Read(versionJSON); err != nil {
		panic(err)
	}
}

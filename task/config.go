package task

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/jwmwalrus/bumpy/internal/config"
	"github.com/urfave/cli/v3"
)

// Config modifies the version config file.
func Config() *cli.Command {
	return &cli.Command{
		Name:            "config",
		Aliases:         []string{"c"},
		Category:        "Control",
		Usage:           "Modify the version config file",
		UsageText:       "config [<flags>...] ...",
		Description:     "Modify the version configuration file and display its contents",
		SkipFlagParsing: false,
		HideHelp:        false,
		Hidden:          false,
		Action:          configAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "persist",
				Usage: "Perform a 'git commit' for the config udate",
			},
			&cli.BoolFlag{
				Name:  "no-fetch",
				Usage: "Do no perform a 'git fetch' operation, persistent as 'config.noFetch'",
			},
			&cli.BoolFlag{
				Name:  "fetch",
				Usage: "Perform a 'git fetch' operation, persistent as 'config.noFetch'",
			},
			&cli.BoolFlag{
				Name:  "no-commit",
				Usage: "Do no perform 'git commit' operations, persistent as 'config.noCommit'",
			},
			&cli.BoolFlag{
				Name:  "commit",
				Usage: "Perform 'git commit' operations, persistent as 'config.noCommit'",
			},
			&cli.StringFlag{
				Name:  "version-prefix",
				Usage: "Subdirectory to store version file, persistent as 'config.VersionPrefix'",
			},
			&cli.StringSliceFlag{
				Name:  "add-npm-prefix",
				Usage: "Add subdirectory to npm prefixes",
			},
			&cli.StringSliceFlag{
				Name:  "remove-npm-prefix",
				Usage: "Remove subdirectory from npm prefixes",
			},
			&cli.BoolFlag{
				Name:  "clear-npm-prefixes",
				Usage: "Clears the list of npm prefixes in the config",
			},
		},
	}
}

func configAction(ctx context.Context, c *cli.Command) (err error) {
	cfg, err := config.Load()
	if err != nil {
		return
	}

	if c.Bool("no-fetch") {
		cfg.NoFetch = true
	} else if c.Bool("fetch") {
		cfg.NoFetch = false
	}

	if c.Bool("no-commit") {
		cfg.NoCommit = true
	} else if c.Bool("commit") {
		cfg.NoCommit = false
	}

	if c.String("version-prefix") != "" {
		cfg.VersionPrefix = c.String("version-prefix")
	}

	if len(c.StringSlice("add-npm-prefix")) > 0 {
		for _, p := range c.StringSlice("add-npm-prefix") {
			cfg.NPMPrefixes = append(cfg.NPMPrefixes, p)
		}
	}

	if len(c.StringSlice("remove-npm-prefix")) > 0 {
		newSlice := []string{}
		// TODO: optimize loop
	outerLoop:
		for _, v := range cfg.NPMPrefixes {
			for _, p := range c.StringSlice("remove-npm-prefix") {
				if v == p {
					continue outerLoop
				}
			}
			newSlice = append(newSlice, v)
		}
		cfg.NPMPrefixes = newSlice
	}

	if c.Bool("clear-npm-prefixes") {
		cfg.NPMPrefixes = []string{}
	}

	if err = cfg.Save(); err != nil {
		return
	}

	bv, err := config.GetBytes()
	if err != nil {
		return
	}

	if c.Bool("persist") {
		if err = cfg.Git.CommitFiles(
			[]string{filepath.Join(".", config.Filename)},
			"Update version config",
		); err != nil {
			return
		}
	}

	fmt.Printf("%v\n", string(bv))
	return
}

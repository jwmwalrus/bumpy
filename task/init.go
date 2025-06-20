package task

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jwmwalrus/bumpy/internal/config"
	"github.com/jwmwalrus/bumpy/version"
	"github.com/urfave/cli/v3"
)

// Init creates an initial version file.
func Init() *cli.Command {
	return &cli.Command{
		Name:            "init",
		Category:        "Control",
		Usage:           "Creates an initial version file",
		UsageText:       "init [--no-fetch]",
		Description:     "Creates an initial configuration file, `.bumpy-ride` at the root of the repository; and creates a version file, `version.file` at the root of the repository or at the location espeficied by the '--version-prefix' flag. The command causes a 'git fetch' as a side-effect, in order to obtain the latest tag",
		SkipFlagParsing: false,
		HideHelp:        false,
		Hidden:          false,
		Action:          initAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "persist",
				Usage: "Perform a 'git commit' for the initialization",
			},
			&cli.BoolFlag{
				Name:  "no-fetch",
				Usage: "Do no perform a 'git fetch' operation, persistent as 'config.noFetch'",
			},
			&cli.BoolFlag{
				Name:  "no-commit",
				Usage: "Do no perform 'git commit' operations, persistent as 'config.noCommit'",
			},
			&cli.StringFlag{
				Name:  "version-prefix",
				Usage: "Subdirectory to store version file, persistent as 'config.VersionPrefix'",
			},
			&cli.StringSliceFlag{
				Name:  "npm-prefix",
				Usage: "ubdirectory to find 'package.json', persistent as 'config.npmPrefixes'",
			},
		},
	}
}

func initAction(ctx context.Context, c *cli.Command) (err error) {
	cfg, configCreated, err := config.LoadOrCreate()
	if err != nil {
		return
	}

	if !configCreated {
		fmt.Printf("Config file already existed!\n")
	}

	versionFile := filepath.Join(cfg.VersionPrefix, version.Filename)
	_, err = os.Stat(versionFile)
	if !os.IsNotExist(err) {
		if !configCreated {
			err = errors.New("Repository is already initialized, isn't it?")
			return
		}
	}

	if c.Bool("no-fetch") {
		if !configCreated {
			fmt.Printf("Overriding `noFetch` in config file")
		}
		cfg.NoFetch = c.Bool("no-fetch")
	}
	if c.Bool("no-commit") {
		if !configCreated {
			fmt.Printf("Overriding `noCommit` in config file")
		}
		cfg.NoCommit = c.Bool("no-commit")
	}
	if c.String("version-prefix") != "" {
		if !configCreated {
			fmt.Printf("Overriding `prefix` in config file")
		}
		cfg.VersionPrefix = c.String("version-prefix")
	}
	if len(c.StringSlice("npm-prefix")) > 0 {
		if !configCreated {
			fmt.Printf("Overriding `npmPrefix` in config file")
		}
		cfg.NPMPrefixes = c.StringSlice("npm-prefix")
	}

	if err = cfg.Save(); err != nil {
		return
	}

	v := version.Version{}
	tag, err := cfg.Git.LatestTag(cfg.NoFetch)
	if err != nil {
		v = version.New()
	} else {
		if err = v.Parse(tag); err != nil {
			return
		}
	}

	if err = v.SaveTo(cfg.VersionPrefix); err != nil {
		return
	}

	if c.Bool("persist") {
		fmt.Printf("\nCommitting files...\n")
		err = cfg.Git.CommitFiles(
			[]string{
				filepath.Join(".", config.Filename),
				versionFile,
			},
			"Init version",
		)
		if err != nil {
			return
		}
	}

	fmt.Printf("Done!\n")
	return
}

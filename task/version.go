package task

import (
	"fmt"

	"github.com/jwmwalrus/bumpy/internal/config"
	"github.com/jwmwalrus/bumpy/version"
	"github.com/urfave/cli/v2"
)

// Version displays the project version
func Version() *cli.Command {
	return &cli.Command{
		Name:            "version",
		Aliases:         []string{"v"},
		Category:        "Informational",
		Usage:           "Display version",
		UsageText:       "version [--short|--long]",
		Description:     "Displays the current version for the repository",
		SkipFlagParsing: false,
		HideHelp:        false,
		Hidden:          false,
		HelpName:        "version",
		Action:          versionAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "short",
				Aliases: []string{"s"},
				Usage:   "Short version",
			},
			&cli.BoolFlag{
				Name:    "long",
				Aliases: []string{"l"},
				Usage:   "Long, detailed version",
			},
		},
	}
}

func versionAction(c *cli.Context) (err error) {
	cfg, err := config.Load()
	if err != nil {
		return
	}

	v := version.Version{}
	if err = v.LoadFrom(cfg.VersionPrefix); err != nil {
		return
	}

	if c.Bool("short") {
		fmt.Printf("%v\n", v.String())
	} else if c.Bool("long") {
		extra := ""
		if v.Pre != "" {
			extra += "\n\tPre: " + v.Pre
		}
		if v.Build != "" {
			extra += "\n\tBuild: " + v.Build
		}

		fmt.Printf("\nVersion: %v\n\tMajor: %v\n\tMinor: %v\n\tPatch: %v%v\n", v.String(), v.Major, v.Minor, v.Patch, extra)
	} else {
		fmt.Printf("\nVersion: %v\n", v.String())
	}

	return
}

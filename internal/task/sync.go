package task

import (
	"fmt"

	"github.com/jwmwalrus/bumpy-ride/internal/git"
	"github.com/jwmwalrus/bumpy-ride/version"
	"github.com/urfave/cli/v2"
)

// Sync synchronizes version file with latest tag
func Sync() *cli.Command {
	return &cli.Command{
		Name:            "sync",
		Category:        "control",
		Usage:           "sync",
		UsageText:       "sync - synchronizes version file",
		Description:     "Synchronizes version file with latest tag",
		SkipFlagParsing: false,
		HideHelp:        false,
		Hidden:          false,
		HelpName:        "init!",
		BashComplete: func(c *cli.Context) {
			// TODO: complete
			fmt.Fprintf(c.App.Writer, "--better\n")
		},
		Action: syncAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "npm-prefix",
				Usage: "Update package.json at the given location",
			},
			&cli.BoolFlag{
				Name:  "no-fetch",
				Usage: "Do no perform a `git fetch` operation",
			},
		},
		OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
			// TODO: complete
			fmt.Fprintf(c.App.Writer, "for shame\n")
			return err
		},
	}
}

func syncAction(c *cli.Context) (err error) {
	v := version.Version{}
	tag := ""
	if tag, err = git.GetLatestTag(c.Bool("no-fetch")); err != nil {
		return
	}

	if err = v.Parse(tag); err != nil {
		return
	}

	if err = v.Save(); err != nil {
		return
	}

	// TODO: update package.json

	fmt.Println("Done!")
	return
}

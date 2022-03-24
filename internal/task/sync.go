package task

import (
	"fmt"

	"github.com/jwmwalrus/bumpy/internal/config"
	"github.com/jwmwalrus/bumpy/pkg/git"
	"github.com/jwmwalrus/bumpy/pkg/version"
	"github.com/urfave/cli/v2"
)

// Sync synchronizes version file with latest tag
func Sync() *cli.Command {
	return &cli.Command{
		Name:            "sync",
		Category:        "Control",
		Usage:           "Synchronizes version file",
		UsageText:       "sync [--npm-prefix PREFIX] [--no-fetch]",
		Description:     "Synchronizes version file with latest tag",
		SkipFlagParsing: false,
		HideHelp:        false,
		Hidden:          false,
		HelpName:        "sync",
		BashComplete: func(c *cli.Context) {
			// TODO: complete
			fmt.Fprintf(c.App.Writer, "--better\n")
		},
		Action: syncAction,
		Flags:  []cli.Flag{},
		OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
			// TODO: complete
			fmt.Fprintf(c.App.Writer, "for shame\n")
			return err
		},
	}
}

func syncAction(c *cli.Context) (err error) {
	var cfg config.Config
	restoreCwd, err := cfg.Load()
	if err != nil {
		return
	}
	defer restoreCwd()

	tag := ""
	if tag, err = git.GetLatestTag(cfg.NoFetch); err != nil {
		return
	}

	v := version.Version{}
	if err = v.Parse(tag); err != nil {
		return
	}

	if err = v.SaveTo(cfg.VersionPrefix); err != nil {
		return
	}

	// TODO: update package.json

	fmt.Printf("Done!\n")
	return
}

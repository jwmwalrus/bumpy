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
		UsageText:       "sync",
		Description:     "Synchronizes version file with latest tag",
		SkipFlagParsing: false,
		HideHelp:        false,
		Hidden:          false,
		HelpName:        "sync",
		Action:          syncAction,
		Flags:           []cli.Flag{},
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

	fmt.Printf("Done!\n")
	return
}

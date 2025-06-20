package task

import (
	"context"
	"fmt"

	"github.com/jwmwalrus/bumpy/internal/config"
	"github.com/jwmwalrus/bumpy/version"
	"github.com/urfave/cli/v3"
)

// Sync synchronizes version file with latest tag.
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
		Action:          syncAction,
		Flags:           []cli.Flag{},
	}
}

func syncAction(ctx context.Context, c *cli.Command) (err error) {
	cfg, err := config.Load()
	if err != nil {
		return
	}

	tag := ""
	if tag, err = cfg.Git.LatestTag(cfg.NoFetch); err != nil {
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

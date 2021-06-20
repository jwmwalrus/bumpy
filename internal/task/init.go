package task

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jwmwalrus/bumpy-ride/internal/git"
	"github.com/jwmwalrus/bumpy-ride/version"
	"github.com/urfave/cli/v2"
)

// Init creates an initial version file
func Init() *cli.Command {
	return &cli.Command{
		Name:            "init",
		Category:        "Control",
		Usage:           "Creates an initial version file",
		UsageText:       "init [--no-fetch]",
		Description:     "Creates an initial version file, using git tags as a hint",
		SkipFlagParsing: false,
		HideHelp:        false,
		Hidden:          false,
		HelpName:        "init",
		BashComplete: func(c *cli.Context) {
			// TODO: complete
			fmt.Fprintf(c.App.Writer, "--better\n")
		},
		Action: initAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-fetch",
				Usage: "Do no perform a 'git fetch' operation",
			},
		},
		OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
			// TODO: complete
			fmt.Fprintf(c.App.Writer, "for shame\n")
			return err
		},
	}
}

func initAction(c *cli.Context) (err error) {
	file := filepath.Join(".", version.VersionFile)
	_, err = os.Stat(file)
	if !os.IsNotExist(err) {
		err = errors.New("Repository is already initialized, isn't it?")
		return
	}

	v := version.Version{}
	tag := ""
	if tag, err = git.GetLatestTag(c.Bool("no-fetch")); err != nil {
		err = v.Save()
		fmt.Println("Done!")
		return
	}

	if err = v.Parse(tag); err != nil {
		return
	}

	if err = v.Save(); err != nil {
		return
	}

	fmt.Println("Done!")
	return
}

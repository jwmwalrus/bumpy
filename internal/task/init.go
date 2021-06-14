package task

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jwmwalrus/bumpy-ride/pkg/version"
	"github.com/urfave/cli/v2"
)

// Init creates an initial version file
func Init() *cli.Command {
	return &cli.Command{
		Name:            "init",
		Category:        "control",
		Usage:           "init",
		UsageText:       "init - creates an initial version file",
		Description:     "Creates an initial version file, using git tags as a hint",
		SkipFlagParsing: false,
		HideHelp:        false,
		Hidden:          false,
		HelpName:        "init!",
		BashComplete: func(c *cli.Context) {
			// TODO: complete
			fmt.Fprintf(c.App.Writer, "--better\n")
		},
		Action: initAction,
		Flags: []cli.Flag{
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

func initAction(c *cli.Context) (err error) {
	file := filepath.Join(".", version.VersionFile)
	_, err = os.Stat(file)
	if !os.IsNotExist(err) {
		err = errors.New("Repository is already initialized, isn't it?")
		return
	}

	v := version.Version{}
	tag := ""
	if tag, err = getLatestTag(c.Bool("no-fetch")); err != nil {
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

func getLatestTag(noFetch bool) (tag string, err error) {
	if !noFetch {
		fmt.Println("Fetching...")
		if _, err = exec.Command("git", "fetch", "--tags").CombinedOutput(); err != nil {
			fmt.Println("...fetching failed!")
		}
	}

	cmd1 := exec.Command("git", "rev-list", "--tags", "--max-count=1")
	output1 := &bytes.Buffer{}
	cmd1.Stdout = output1
	if err = cmd1.Run(); err != nil {
		return
	}
	hash := string(output1.Bytes())
	hash = strings.TrimSuffix(hash, "\n")

	cmd2 := exec.Command("git", "describe", "--tags", hash)
	output2 := &bytes.Buffer{}
	cmd2.Stdout = output2
	if err = cmd2.Run(); err != nil {
		return
	}

	tag = string(output2.Bytes())
	tag = strings.TrimSuffix(tag, "\n")

	return
}

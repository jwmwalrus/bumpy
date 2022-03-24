package task

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jwmwalrus/bumpy/internal/config"
	"github.com/jwmwalrus/bumpy/pkg/git"
	"github.com/jwmwalrus/bumpy/pkg/version"
	"github.com/russross/blackfriday/v2"
	"github.com/urfave/cli/v2"
)

// Tag commits the ChangeLog and adds a tag to it
func Tag() *cli.Command {
	return &cli.Command{
		Name:            "tag",
		Aliases:         []string{"t"},
		Category:        "Git",
		Usage:           "Tags the ChangeLog",
		UsageText:       "tag [--changelog-name NAME]",
		Description:     "Commits ChangeLog.md and tags the commit with the latest version",
		SkipFlagParsing: false,
		HideHelp:        false,
		Hidden:          false,
		HelpName:        "tag",
		BashComplete: func(c *cli.Context) {
			// TODO: complete
			fmt.Fprintf(c.App.Writer, "--better\n")
		},
		Action: tagAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "changelog-name",
				Usage: "Name (including extension) of the ChangeLog file",
			},
		},
		OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
			// TODO: complete
			fmt.Fprintf(c.App.Writer, "for shame\n")
			return err
		},
	}
}

func tagAction(c *cli.Context) (err error) {
	var cfg config.Config
	restoreCwd, err := cfg.Load()
	if err != nil {
		return
	}
	defer restoreCwd()

	fmt.Printf("\nLoading current version file...\n")
	v := version.Version{}
	if err = v.LoadFrom(cfg.VersionPrefix); err != nil {
		return
	}

	fmt.Printf("\tVersion to use as tag: %v\n", v.String())
	filename := c.String("changelog-name")

	if filename, err = resolveChangeLogFilename(filename); err != nil {
		return
	}

	fmt.Printf("\nLoading %v...\n", filename)

	msg := getChangeLogMessage(v, filename)

	if err = git.CommitFiles([]string{filename}, "Update ChangeLog"); err != nil {
		return
	}

	if err = git.CreateTag(v.String(), strings.TrimSuffix(msg, "\n")); err != nil {
		return
	}

	fmt.Printf("\nDone!\n")

	return
}

func getChangeLogMessage(v version.Version, filename string) (msg string) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	var bv []byte
	bv, err = ioutil.ReadAll(file)
	if err != nil {
		return
	}

	fmt.Printf("\nParsing %v...\n", filename)
	md := blackfriday.New()

	var node *blackfriday.Node
	if node = md.Parse(bv); node == nil {
		err = errors.New("Error parsing Markdown")
		return
	}

	node.Walk(func(n *blackfriday.Node, e bool) blackfriday.WalkStatus {
		if e && n.Type == blackfriday.Heading && n.FirstChild != nil && n.Next != nil && n.Next.Type == blackfriday.Paragraph {
			if strings.Contains(string(n.FirstChild.Literal), v.StringNoV()) {
				msg = string(n.Next.FirstChild.Literal)
				return blackfriday.Terminate
			}
		}
		return blackfriday.GoToNext
	})

	if msg == "" {
		msg = "New version"
	}

	return
}

func resolveChangeLogFilename(filename string) (string, error) {
	if filename == "" {
		fmt.Printf("\nLooking for a ChangeLog file\n")
		commonNames := []string{
			"CHANGELOG.md",
			"ChangeLog.md",
			"Changelog.md",
			"changelog.md",
			"HISTORY.md",
			"History.md",
			"history.md",
			"NEWS.md",
			"News.md",
			"news.md",
			"RELEASES.md",
			"Releases.md",
			"releases.md",
		}

		for _, fn := range commonNames {
			_, err := os.Stat(fn)
			if os.IsNotExist(err) {
				continue
			}
			fmt.Printf("\tFound filename: %v\n", fn)
			filename = fn
			break
		}

		if filename == "" {
			return filename, fmt.Errorf("\tUnable to find ChangeLog")
		}
	}
	return filename, nil
}

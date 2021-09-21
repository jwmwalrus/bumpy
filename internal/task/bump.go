package task

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/jwmwalrus/bumpy-ride/internal/git"
	"github.com/jwmwalrus/bumpy-ride/version"
	"github.com/urfave/cli/v2"
)

// Bump bumps version
func Bump() *cli.Command {
	return &cli.Command{
		Name:            "bump",
		Aliases:         []string{"b"},
		Category:        "Git",
		Usage:           "Increase current version",
		UsageText:       "bump [--major|--minor|--patch] [--pre PRE] [--build BUILD] ...",
		Description:     "Increases the current version according to the given options",
		SkipFlagParsing: false,
		HideHelp:        false,
		Hidden:          false,
		HelpName:        "bump",
		BashComplete: func(c *cli.Context) {
			// TODO: complete
			fmt.Fprintf(c.App.Writer, "--better\n")
		},
		Before: checkVersionInSync,
		Action: bumpAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "major",
				Aliases: []string{"maj"},
				Usage:   "Increase major version number",
			},
			&cli.BoolFlag{
				Name:    "minor",
				Aliases: []string{"min"},
				Usage:   "Increase minor version number",
			},
			&cli.BoolFlag{
				Name:    "patch",
				Aliases: []string{"p"},
				Usage:   "Increase patch version number",
			},
			&cli.StringFlag{
				Name:  "pre",
				Usage: "Assign `PRE` to the prerelease version string",
			},
			&cli.StringFlag{
				Name:  "build",
				Usage: "Assign `BUILD` to the build version string",
			},
			&cli.StringFlag{
				Name:  "npm-prefix",
				Usage: "Update package.json at the given `PREFIX`  location",
			},
			&cli.BoolFlag{
				Name:  "no-fetch",
				Usage: "Do no perform a 'git fetch' operation",
			},
			&cli.BoolFlag{
				Name:  "no-commit",
				Usage: "Do no commit version file(s)",
			},
		},
		OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
			// TODO: complete
			fmt.Fprintf(c.App.Writer, "for shame\n")
			return err
		},
	}
}

func bumpAction(c *cli.Context) (err error) {
	var v version.Version

	if err = v.Load(); err != nil {
		return
	}

	rest := c.Args().Slice()

	if c.Bool("major") {
		v.Major = v.Major + 1
		v.Minor = 0
		v.Patch = 0
		v.Pre = ""
		v.Build = ""
	} else if c.Bool("minor") {
		v.Minor = v.Minor + 1
		v.Patch = 0
		v.Pre = ""
		v.Build = ""
	} else if c.Bool("patch") {
		v.Patch = v.Patch + 1
		v.Pre = ""
		v.Build = ""
	} else {
		if len(rest) > 1 {
			err = errors.New("Too many options provided")

		} else if len(rest) == 1 {
			if err = v.Parse(rest[0]); err != nil {
				return
			}
		}
	}

	if c.String("pre") != "" {
		v.Pre = c.String("pre")
	}
	if c.String("build") != "" {
		v.Build = c.String("build")
	}

	if err = v.Save(); err != nil {
		return
	}

	sList := []string{filepath.Join("", version.VersionFile)}

	if c.String("npm-prefix") != "" {
		var jsonFiles []string
		if jsonFiles, err = updatePackageJSON(c.String("npm-prefix"), v); err != nil {
			return
		}

		for _, f := range jsonFiles {
			sList = append(sList, f)
		}
	}

	if !c.Bool("no-commit") {
		if err = git.CommitFiles(sList, "Bump version"); err != nil {
			return
		}
	}

	fmt.Printf("Done!\n")
	return
}

func checkVersionInSync(c *cli.Context) (err error) {
	var vFromFile version.Version
	var tag string

	if err = vFromFile.Load(); err != nil {
		return
	}

	if tag, err = git.GetLatestTag(c.Bool("no-fetch")); err != nil {
		fmt.Printf("WARNING, unable to obtain latest tag: %v\n", err)
		err = nil
		return
	}

	var ok bool
	if ok, err = vFromFile.EqualsString(tag); err != nil || !ok {
		return
	}

	if !ok {
		err = errors.New("Version in file does not match latest tag. Please sync")
	}

	return
}

func updatePackageJSON(prefix string, v version.Version) (files []string, err error) {
	if _, err = exec.Command("npm", "version", "--prefix", prefix, "--no-git-tag-version", v.String()).CombinedOutput(); err != nil {
		return
	}
	files = append(files, filepath.Join(prefix, "package.json"))
	files = append(files, filepath.Join(prefix, "package-lock.json"))
	return
}

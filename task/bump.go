package task

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/jwmwalrus/bnp/git"
	"github.com/jwmwalrus/bumpy/internal/config"
	"github.com/jwmwalrus/bumpy/version"
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
		Before:          checkVersionInSync,
		Action:          bumpAction,
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
		},
	}
}

func bumpAction(c *cli.Context) (err error) {
	var cfg config.Config
	restoreCwd, err := cfg.Load()
	if err != nil {
		return
	}
	defer restoreCwd()

	var v version.Version
	if err = v.LoadFrom(cfg.VersionPrefix); err != nil {
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

	if err = v.SaveTo(cfg.VersionPrefix); err != nil {
		return
	}

	slist := []string{
		filepath.Join(".", config.Filename),
		filepath.Join(cfg.VersionPrefix, version.Filename),
	}

	for _, p := range cfg.NPMPrefixes {
		var jsonFiles []string
		if jsonFiles, err = updatePackageJSON(p, v); err != nil {
			return
		}

		for _, f := range jsonFiles {
			slist = append(slist, f)
		}
	}

	if !cfg.NoCommit {
		if err = git.CommitFiles(slist, "Bump version"); err != nil {
			return
		}
	}

	fmt.Printf("Done!\n\nNext tag will be: %v\n", v.String())
	return
}

func checkVersionInSync(c *cli.Context) (err error) {
	var cfg config.Config
	restoreCwd, err := cfg.Load()
	if err != nil {
		return
	}
	defer restoreCwd()

	var vFromFile version.Version
	var tag string

	if err = vFromFile.LoadFrom(cfg.VersionPrefix); err != nil {
		return
	}

	if tag, err = git.GetLatestTag(cfg.NoFetch); err != nil {
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

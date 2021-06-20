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
		Category:        "control",
		Usage:           "bump [--major|--minor|--patch] [--pre PRE] [--build BUILD] ...",
		UsageText:       "bump - increase current version",
		Description:     "Increases the current version according to the given flaGS",
		SkipFlagParsing: false,
		HideHelp:        false,
		Hidden:          false,
		HelpName:        "bump!",
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

	if err = git.CommitFiles(sList, "Bump version"); err != nil {
		return
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

/*
func bumpAction(c *cli.Context) (err error) {
		const path = "/playback/"

		uri := base.Conf.Server.GetURL() + path

		res, err := http.Get(uri)
		if err != nil {
			log.Error(err)
			return
		}
		defer res.Body.Close()

		r, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Error(err)
		}

		// TODO: prettify

		if c.Bool("json") {
			var out bytes.Buffer
			json.Indent(&out, r, "", "  ")
			out.WriteTo(os.Stdout)
		} else if c.Bool("raw") {
			fmt.Println(string(r))
		} else {
			type jt struct {
				TrackID  int64
				Location string
			}
			j := jt{}
			err = json.Unmarshal(r, &j)

			tbl := table.New("ID", "Location")
			tbl.AddRow(j.TrackID, j.Location)
			tbl.Print()
		}

	return
}

func playbackPlayAction(c *cli.Context) (err error) {
	const path = "/playback/"
	body := base.PlaybackReqJSON{}

	rest := c.Args().Slice()
	if len(rest) == 1 && bnp.IsJSON(rest[0]) {
		_ = json.Unmarshal([]byte(rest[0]), &body)
	} else {
		for _, v := range rest {
			var u *url.URL
			if u, err = url.Parse(v); err != nil {
				return
			}
			if u.Scheme == "" {
				u.Scheme = "file"
			}
			body.Locations = append(body.Locations, u.String())
		}
	}
	body.Action = base.PlaybackReqActionPlay
	body.Force = c.Bool("force")

	uri := base.Conf.Server.GetURL() + path

	jm, _ := json.Marshal(&body)
	data := bytes.NewBuffer(jm)

	res, err := http.Post(uri, "application/json", data)
	if err != nil {
		log.Error(err)
		return
	}
	defer res.Body.Close()

	r, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(err)
	}

	fmt.Println(string(r))
	return
}

func playbackPauseAction(c *cli.Context) (err error) {
	const path = "/playback/"
	body := base.PlaybackReqJSON{}

	rest := c.Args().Slice()
	if len(rest) > 0 {
		err = errors.New("Too many values in command")
		return
	}
	body.Action = base.PlaybackReqActionPause

	uri := base.Conf.Server.GetURL() + path

	jm, _ := json.Marshal(&body)
	data := bytes.NewBuffer(jm)

	res, err := http.Post(uri, "application/json", data)
	if err != nil {
		log.Error(err)
		return
	}
	defer res.Body.Close()

	r, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(err)
	}

	fmt.Println(string(r))
	return
}

func playbackStopNextPreviousAction(c *cli.Context) (err error) {
	// TODO: implement
	fmt.Println("TODO")
	return
}

func playbackJumpAction(c *cli.Context) (err error) {
	// TODO: implement
	fmt.Println("TODO")
	return
}
*/

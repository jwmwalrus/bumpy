package version

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

const (
	// Filename names the version file
	Filename = "version.json"
)

// Version Basic SemVer structure
type Version struct {
	Major int    `json:"major"`
	Minor int    `json:"minor"`
	Patch int    `json:"patch"`
	Pre   string `json:"pre"`
	Build string `json:"build"`

	revision string
	time     time.Time
	modified bool
}

// New returns an initial version
func New() (v Version) {
	v.Minor = 1
	return
}

// Equals checks if two versions are identical
func (v *Version) Equals(r Version) bool {
	return v.Major == r.Major &&
		v.Minor == r.Minor &&
		v.Patch == r.Patch &&
		v.Pre == r.Pre &&
		v.Build == r.Build
}

// EqualsString checks if version is identical to string
func (v *Version) EqualsString(s string) (ok bool, err error) {
	var r Version
	if err = r.Parse(s); err != nil {
		return
	}

	ok = v.Equals(r)
	return
}

// GetLastCommit returns the embedded vcs-related information
func (v *Version) GetLastCommit() (string, time.Time, bool) {
	return v.revision, v.time, v.modified
}

// Load loads the version file from the current working directory
func (v *Version) Load() (err error) {
	err = v.LoadFrom(".")
	return
}

// LoadFrom loads the version file from the given directory
func (v *Version) LoadFrom(dir string) (err error) {
	file := filepath.Join(dir, Filename)
	_, err = os.Stat(file)
	if os.IsNotExist(err) {
		err = fmt.Errorf("The given path does not exist: %v", file)
		return
	}

	jsonFile, err := os.Open(file)
	if err != nil {
		return
	}
	defer jsonFile.Close()

	var byteValue []byte
	byteValue, err = ioutil.ReadAll(jsonFile)
	if err != nil {
		return
	}

	if err = v.Read(byteValue); err != nil {
		return
	}

	err = v.getBuildInfo()
	return
}

// Parse parses a version string into its fields
func (v *Version) Parse(s string) error {
	var core, pre, build []byte
	var inCore, inPre, inBuild bool
	inCore = true
	for i, c := range s {
		if i == 0 && c == 'v' {
			continue
		} else if !inPre && c == '-' {
			inCore = false
			inPre = true
		} else if !inBuild && c == '+' {
			inCore = false
			inPre = false
			inBuild = true
		} else {
			if inCore {
				core = append(core, byte(c))
			} else if inPre {
				pre = append(pre, byte(c))
			} else if inBuild {
				build = append(build, byte(c))
			} else {
				return fmt.Errorf("Error parsing version string at character '%v', position %v",
					c, i+1)
			}
		}
	}

	a := strings.Split(string(core), ".")

	if len(a) != 3 {
		return fmt.Errorf("Version string does not follow a major.minor.patch pattern")
	}

	mmp := make([]int64, 3)
	for i, x := range a {
		var err error
		if mmp[i], err = strconv.ParseInt(x, 10, 32); err != nil {
			return fmt.Errorf("Cover version #%v is not an integer", i+1)
		}
	}
	v.Major = int(mmp[0])
	v.Minor = int(mmp[1])
	v.Patch = int(mmp[2])
	v.Pre = string(pre)
	v.Build = string(build)

	return nil
}

// Read reads the version from the given bytes
func (v *Version) Read(b []byte) (err error) {
	err = json.Unmarshal(b, v)
	return
}

// Save saves the version file to the current working directory
func (v *Version) Save() (err error) {
	err = v.SaveTo(".")
	return
}

// SaveTo saves the version file to the given directory
func (v *Version) SaveTo(dir string) (err error) {
	var file []byte
	file, err = json.Marshal(v)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(filepath.Join(dir, Filename), file, 0644)
	return
}

func (v *Version) String() string {
	return "v" + v.StringNoV()
}

// StringNoV returns the version string, without a "v" prefix
func (v *Version) StringNoV() (out string) {
	out = strconv.Itoa(v.Major) + "." + strconv.Itoa(v.Minor) + "." +
		strconv.Itoa(v.Patch)

	if v.Pre != "" {
		out += "-" + v.Pre
	}

	if v.Build != "" {
		out += "+" + v.Build
	}

	return
}

func (v *Version) getBuildInfo() (err error) {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		err = fmt.Errorf("Error obtaining runtime debug.BuildInfo")
		return
	}

	var vRevision string
	var vTime time.Time
	var isGit, vModified bool
	for _, s := range bi.Settings {
		switch s.Key {
		case "vcs":
			isGit = s.Value == "git"
		case "vcs.revision":
			vRevision = s.Value
		case "vcs.time":
			vTime, _ = time.Parse(time.RFC3339, s.Value)
		case "vcs.modified":
			vModified = s.Value == "true"
		}
	}

	if !isGit {
		err = fmt.Errorf("Unsupported VCS")
		return
	}

	v.revision = vRevision
	v.time = vTime
	v.modified = vModified

	return
}

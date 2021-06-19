package version

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	// VersionFile names the version file
	VersionFile = "version.json"
)

// Version Basic SemVer structure
type Version struct {
	Major int    `json:"major"`
	Minor int    `json:"minor"`
	Patch int    `json:"patch"`
	Pre   string `json:"pre"`
	Build string `json:"build"`
}

// Load loads the version file from the current working directory
func (v *Version) Load() (err error) {
	err = v.LoadFrom(".")
	return
}

// LoadFrom loads the version file from the given path
func (v *Version) LoadFrom(path string) (err error) {
	file := filepath.Join(path, VersionFile)
	_, err = os.Stat(file)
	if os.IsNotExist(err) {
		err = errors.New("The given path does not exist")
		return
	}

	// var jsonFile *os.File
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

	err = v.Read(byteValue)
	return
}

// NoPrefix returns the version string, without a "v" prefix
func (v *Version) NoPrefix() (out string) {
	out = strconv.Itoa(v.Major) + "." + strconv.Itoa(v.Minor) + "." + strconv.Itoa(v.Patch)

	if v.Pre != "" {
		out += "-" + v.Pre
	}

	if v.Build != "" {
		out += "+" + v.Build
	}

	return
}

// Parse parses a version string into its fields
func (v *Version) Parse(s string) (err error) {
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
				err = fmt.Errorf("Error parsing version string at character '%v', position %v", c, i+1)
				return
			}
		}
	}

	a := strings.Split(string(core), ".")

	if len(a) != 3 {
		err = errors.New("Version string does not follow a major.minor.patch pattern")
		return
	}

	mmp := make([]int64, 3)
	for i, x := range a {
		if mmp[i], err = strconv.ParseInt(x, 10, 32); err != nil {
			err = fmt.Errorf("Cover version #%v is not an integer", i+1)
			return
		}
	}
	v.Major = int(mmp[0])
	v.Minor = int(mmp[1])
	v.Patch = int(mmp[2])
	v.Pre = string(pre)
	v.Build = string(build)

	return
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

// SaveTo saves the version file to the given path
func (v *Version) SaveTo(path string) (err error) {
	var file []byte
	file, err = json.Marshal(v)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(filepath.Join(path, VersionFile), file, 0644)
	return
}

// String returns the version string
func (v *Version) String() (out string) {
	out = "v" + v.NoPrefix()
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

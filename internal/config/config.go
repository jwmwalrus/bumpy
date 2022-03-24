package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jwmwalrus/bumpy/pkg/git"
)

const (
	// Filename names the config file
	Filename = ".bumpy-ride"
)

// Config defines the bumpy-ride configuration file
type Config struct {
	NoFetch       bool     `json:"noFetch"`
	NoCommit      bool     `json:"noCommit"`
	VersionPrefix string   `json:"versionPrefix"`
	NPMPrefixes   []string `json:"npmPrefixes"`
}

// New returns an initial Config
func New() (cfg Config) {
	cfg.VersionPrefix = "."
	cfg.NPMPrefixes = []string{}
	return
}

// Load loads the configuration file, which must exist
func (cfg *Config) Load() (fn git.RestoreCwdFunc, err error) {
	fn, err = git.MoveToRootDir()
	if err != nil {
		return
	}

	if _, err = os.Stat(filepath.Join(".", Filename)); os.IsNotExist(err) {
		return
	}

	err = cfg.Read()
	return
}

// LoadOrCreate loads the configuration file if it exists, or creates it otherwise
func (cfg *Config) LoadOrCreate() (created bool, fn git.RestoreCwdFunc, err error) {
	fn, err = git.MoveToRootDir()
	if err != nil {
		return
	}

	if _, err = os.Stat(filepath.Join(".", Filename)); os.IsNotExist(err) {
		if cfg == nil {
			cfg = &Config{}
		}
		*cfg = New()
		if err = cfg.Save(); err != nil {
			return
		}
		created = true
		return
	}

	err = cfg.Read()
	return
}

// Read reads the configuration file
func (cfg *Config) Read() (err error) {
	bv, err := GetBytes()
	if err != nil {
		return
	}

	err = json.Unmarshal(bv, cfg)
	return
}

// Save writes the configuration file
func (cfg *Config) Save() (err error) {
	var bv []byte
	bv, err = json.MarshalIndent(*cfg, "", "  ")
	if err != nil {
		return
	}

	err = ioutil.WriteFile(filepath.Join(".", Filename), bv, 0644)
	return
}

// GetBytes reads the configuration file and returns its bytes
func GetBytes() (bv []byte, err error) {
	jsonFile, err := os.Open(filepath.Join(".", Filename))
	if err != nil {
		return
	}
	defer jsonFile.Close()

	bv, err = ioutil.ReadAll(jsonFile)
	return
}

package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jwmwalrus/bnp/git"
	"github.com/jwmwalrus/bnp/onerror"
)

const (
	// Filename names the config file
	Filename = ".bumpy-ride"
)

// Config defines the bumpy-ride configuration file
type Config struct {
	NoFetch       bool          `json:"noFetch"`
	NoCommit      bool          `json:"noCommit"`
	VersionPrefix string        `json:"versionPrefix"`
	NPMPrefixes   []string      `json:"npmPrefixes"`
	Git           git.Interface `json:"-"`
}

// New returns an initial Config
func New() *Config {
	cfg := &Config{}
	cfg.gitLoad()

	cfg.VersionPrefix = "."
	cfg.NPMPrefixes = []string{}

	return cfg
}

// Load loads the configuration file, which must exist
func Load() (cfg *Config, err error) {
	cfg = &Config{}

	if _, err = os.Stat(filepath.Join(".", Filename)); os.IsNotExist(err) {
		cfg = nil
		return
	}

	if err = cfg.Read(); err != nil {
		cfg = nil
		return
	}

	cfg.gitLoad()
	return
}

// LoadOrCreate loads the configuration file if it exists, or creates it otherwise
func LoadOrCreate() (cfg *Config, created bool, err error) {
	cfg = &Config{}

	if _, err = os.Stat(filepath.Join(".", Filename)); errors.Is(err, os.ErrNotExist) {
		cfg = New()
		if err = cfg.Save(); err != nil {
			cfg = nil
			return
		}
		created = true
		return
	}

	if err = cfg.Read(); err != nil {
		cfg = nil
		return
	}

	cfg.gitLoad()
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
	bv, err := json.MarshalIndent(*cfg, "", "  ")
	if err != nil {
		return
	}

	err = ioutil.WriteFile(filepath.Join(".", Filename), bv, 0644)
	return
}

func (cfg *Config) gitLoad() {
	cwd, _ := os.Getwd()
	var err error
	cfg.Git, err = git.NewInterface(cwd)
	onerror.Fatal(err)
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

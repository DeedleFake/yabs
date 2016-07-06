package main

import (
	"io"

	"github.com/naoina/toml"
)

// A Config maps to the TOML layout of a timeline config file.
type Config struct {
	Source string `toml:"source"`
	Dest   string `toml:"dest"`

	Regular struct {
		NameScheme string `toml:"namescheme"`
		Writable   bool   `toml:"writable"`
	} `toml:"regular"`
}

// LoadConfig loads a new Config from r.
func LoadConfig(r io.Reader) (*Config, error) {
	d := toml.NewDecoder(r)

	var cfg Config
	err := d.Decode(&cfg)
	return &cfg, err
}

// A Timeline stores information about the current state of a
// timeline. It's methods and fields allow operations on a timeline,
// such as creating new snapshots, finding existing ones, and deleting
// them.
type Timeline struct {
}

// Timeline loads a Timeline from a Config. This is the main way of
// obtaining a valid Timeline.
func (cfg *Config) Timeline() (*Timeline, error) {
	panic("Not implemented.")
}

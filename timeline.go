package main

import (
	"io"

	"github.com/naoina/toml"
)

type Config struct {
	Source string `toml:"source"`
	Dest   string `toml:"dest"`

	Regular struct {
		NameScheme string `toml:"namescheme"`
	} `toml:"regular"`
}

func LoadConfig(r io.Reader) (*Config, error) {
	d := toml.NewDecoder(r)

	var cfg Config
	err := d.Decode(&cfg)
	return &cfg, err
}

type Timeline struct {
}

func (cfg *Config) Timeline() (*Timeline, error) {
	panic("Not implemented.")
}

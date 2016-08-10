package main

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/naoina/toml"
)

// A Config maps to the TOML layout of a timeline config file.
type Config struct {
	Source string `toml:"source"`
	Dest   string `toml:"dest"`

	NameScheme string `toml:"namescheme"`
	UTC        bool   `toml:"utc"`
	Writable   bool   `toml:"writable"`

	Allowed ConfigAllowed `toml:"allowed"`
}

type ConfigAllowed struct {
	Age string `toml:"age"`
	Num string `toml:"num"`
}

// LoadConfig loads a new Config from r.
func LoadConfig(r io.Reader) (*Config, error) {
	d := toml.NewDecoder(r)

	cfg := &Config{
		NameScheme: "Stamp",

		Allowed: ConfigAllowed{
			Age: "1000h",
			Num: "100",
		},
	}
	err := d.Decode(cfg)
	return cfg, err
}

func (cfg *Config) Update(ctx context.Context) error {
	if cfg.Source == "" {
		return errors.New("Config has no source specified.")
	}
	if cfg.Dest == "" {
		return errors.New("Config has no destination specified.")
	}

	panic("Not implemented.")
}

var timeNames = map[string]string{
	"ANSIC":       time.ANSIC,
	"UnixDate":    time.UnixDate,
	"RubyDate":    time.RubyDate,
	"RFC822":      time.RFC822,
	"RFC822Z":     time.RFC822Z,
	"RFC850":      time.RFC850,
	"RFC1123":     time.RFC1123,
	"RFC1123Z":    time.RFC1123Z,
	"RFC3339":     time.RFC3339,
	"RFC3339Nano": time.RFC3339Nano,
	"Kitchen":     time.Kitchen,
	"Stamp":       time.Stamp,
	"StampMilli":  time.StampMilli,
	"StampMicro":  time.StampMicro,
	"StampNano":   time.StampNano,
}

func FormatTime(t time.Time, f string) string {
	if n, ok := timeNames[f]; ok {
		f = n
	}

	return t.Format(f)
}

package main

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/naoina/toml"
)

// A Config maps to the TOML layout of a timeline config file.
type Config struct {
	Source string `toml:"source"`
	Dest   string `toml:"dest"`

	NameFormat string `toml:"nameformat"`
	UTC        bool   `toml:"utc"`
	Writable   bool   `toml:"writable"`

	Allowed ConfigAllowed `toml:"allowed"`
}

type ConfigAllowed struct {
	Num int    `toml:"num"`
	Age string `toml:"age"`
}

// LoadConfig loads a new Config from r.
func LoadConfig(r io.Reader) (*Config, error) {
	d := toml.NewDecoder(r)

	cfg := &Config{
		NameFormat: "Stamp",

		Allowed: ConfigAllowed{
			Num: -1,
			Age: "3000h",
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

	now := time.Now()

	err := CreateSnapshot(ctx,
		cfg.Source,
		filepath.Join(cfg.Dest, now.Format(TimeFormat(cfg.NameFormat))),
		cfg.Writable,
	)
	if err != nil {
		return err
	}

	err = cfg.deleteByNum(ctx)
	if err != nil {
		return err
	}

	//err = cfg.deleteByAge(ctx, now)
	//if err != nil {
	//	return err
	//}

	return nil
}

func (cfg *Config) delete(ctx context.Context, del []os.FileInfo) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, del := range del {
		func(del string) {
			eg.Go(func() error {
				return DeleteSubvol(ctx, del)
			})
		}(filepath.Join(cfg.Dest, del.Name()))
	}
	return eg.Wait()
}

func (cfg *Config) deleteByNum(ctx context.Context) error {
	if cfg.Allowed.Num <= 0 {
		return nil
	}

	dir, err := os.Open(cfg.Dest)
	if err != nil {
		return err
	}
	defer dir.Close()

	c, err := dir.Readdir(-1)
	if err != nil {
		return err
	}
	if len(c) < cfg.Allowed.Num {
		return nil
	}
	sort.Sort(FileInfoByTimestamp{
		fi: c,
		f:  TimeFormat(cfg.NameFormat),
	})

	return cfg.delete(ctx, c[cfg.Allowed.Num:])
}

func (cfg *Config) deleteByAge(ctx context.Context, now time.Time) error {
	age, err := time.ParseDuration(cfg.Allowed.Age)
	if err != nil {
		// Skip deleting by age if duration is invalid.
		return nil
	}

	dir, err := os.Open(cfg.Dest)
	if err != nil {
		return err
	}
	defer dir.Close()

	c, err := dir.Readdir(-1)
	if err != nil {
		return err
	}
	sort.Sort(FileInfoByTimestamp{
		fi: c,
		f:  TimeFormat(cfg.NameFormat),
	})

	newest := now.Add(-age)
	first := sort.Search(len(c), func(i int) bool {
		return c[i].ModTime().Before(newest)
	})
	if first == len(c) {
		return nil
	}

	return cfg.delete(ctx, c[first:])
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

func TimeFormat(f string) string {
	if n, ok := timeNames[f]; ok {
		return n
	}

	return f
}

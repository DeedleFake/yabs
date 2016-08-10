package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"sync"

	"golang.org/x/sys/unix"
)

var L = struct {
	I *log.Logger
	E *log.Logger
}{
	I: log.New(os.Stdout, "[INFO] ", log.LstdFlags),
	E: log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
}

func timelines(root string) ([]string, error) {
	dir, err := os.Open(root)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	fi, err := dir.Stat()
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("%q is not a directory", root)
	}

	c, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	sort.Sort(FileInfoByName(c))

	list := make([]string, 0, len(c))
	for _, entry := range c {
		if !entry.Mode().IsRegular() {
			continue
		}

		list = append(list, entry.Name())
	}

	return list, nil
}

func update(ctx context.Context, cpath string) error {
	file, err := os.Open(cpath)
	if err != nil {
		return err
	}
	defer file.Close()

	cfg, err := LoadConfig(file)
	if err != nil {
		return err
	}

	return cfg.Update(ctx)
}

func main() {
	u, err := user.Current()
	if err != nil {
		L.E.Fatalf("Failed to get current user: %v", err)
	}
	if u.Uid != "0" {
		L.E.Fatalf("%v must be run as root.", os.Args[0])
	}

	var flags struct {
		configRoot string
	}
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v [options] <timeline>\n", os.Args[0])
		fmt.Fprintln(os.Stderr)

		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr)

		fmt.Fprintln(os.Stderr, "Pseudo timelines:")
		fmt.Fprintln(os.Stderr, "  list-timelines")
		fmt.Fprintln(os.Stderr, "      List all timeline configs.")
		fmt.Fprintln(os.Stderr, "  update-all")
		fmt.Fprintln(os.Stderr, "      Update all timelines.")
	}
	flag.StringVar(&flags.configRoot, "confdir", "/etc/yabs/", "The directory that the timeline configs are in.")
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	ctx := SignalContext(context.Background(),
		os.Interrupt,
		unix.SIGTERM,
	)

	switch timeline := flag.Arg(0); timeline {
	case "list-timelines":
		tl, err := timelines(flags.configRoot)
		if err != nil {
			L.E.Fatalf("Failed to get list of timelines: %v", err)
		}

		for _, tl := range tl {
			fmt.Println(tl)
		}

	case "update-all":
		tl, err := timelines(flags.configRoot)
		if err != nil {
			L.E.Fatalf("Failed to get list of timelines: %v", err)
		}

		var wg sync.WaitGroup
		for _, tl := range tl {
			wg.Add(1)
			go func(tl string) {
				defer wg.Done()

				err := update(ctx, filepath.Join(flags.configRoot, tl))
				if err != nil {
					L.E.Printf("Failed to update %q: %v", tl, err)
					return
				}

				L.I.Printf("Updated %q.", tl)
			}(tl)
		}
		wg.Wait()

	default:
		err := update(ctx, filepath.Join(flags.configRoot, timeline))
		if err != nil {
			L.E.Fatalf("Failed to update %q: %v", timeline, err)
			return
		}

		L.I.Printf("Updated %q.", timeline)
	}
}

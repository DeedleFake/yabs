package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
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
	panic("Not implemented.")
}

func main() {
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

	ctx := SignalContext(context.Background(), os.Interrupt)

	switch timeline := flag.Arg(0); timeline {
	case "list-timelines":
		tl, err := timelines(flags.configRoot)
		if err != nil {
			L.E.Printf("Failed to get list of timelines: %v", err)
			os.Exit(1)
		}

		for _, tl := range tl {
			fmt.Println(tl)
		}

	case "update-all":
		tl, err := timelines(flags.configRoot)
		if err != nil {
			L.E.Printf("Failed to get list of timelines: %v", err)
			os.Exit(1)
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
			L.E.Printf("Failed to update %q: %v", timeline, err)
		}

		L.I.Printf("Update %q.", timeline)
	}
}

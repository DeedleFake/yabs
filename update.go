package main

import (
	"flag"
	"fmt"
	"os"
)

type Update struct {
	root string
}

func (u *Update) Usage(flagSet *flag.FlagSet) {
	fmt.Fprintf(os.Stderr, "Usage: %v [global options] update <timeline config name>\n", os.Args[0])
	fmt.Fprintln(os.Stderr)

	fmt.Fprintln(os.Stderr, `The update command updates a timeline, creating and deleting snapshots
as necessary based on that timeline's configuration.`)
}

func (u *Update) Main(args []string) {
	flagSet := flag.NewFlagSet(args[0], flag.ExitOnError)
	flagSet.Usage = func() { u.Usage(flagSet) }
	flagSet.Parse(args[1:])

	if flagSet.NArg() == 0 {
		flagSet.Usage()
		os.Exit(2)
	}
}

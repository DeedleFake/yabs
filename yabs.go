package main

import (
	"flag"
	"fmt"
	"os"
)

func globalUsage(flagSet *flag.FlagSet) {
	fmt.Fprintf(os.Stderr, "Usage: %v [global options] <command> [command options]\n", os.Args[0])
	fmt.Fprintln(os.Stderr)

	fmt.Fprintln(os.Stderr, "Global options:")
	flagSet.PrintDefaults()
	fmt.Fprintln(os.Stderr)

	fmt.Fprintln(os.Stderr, `Commands:
  update <timeline config name>
        Updates a snapshot timeline.

Follow a command with --help for more detailed information.`)
}

func main() {
	var flags struct {
		cRoot string
	}

	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagSet.Usage = func() { globalUsage(flagSet) }
	flagSet.StringVar(&flags.cRoot, "c", "/etc/yabs/timeline.d/", "The dir to look for configs in.")
	flagSet.Parse(os.Args[1:])

	switch cmd := flagSet.Arg(0); cmd {
	case "update":
		(&Update{
			root: flags.cRoot,
		}).Main(flagSet.Args())

	case "help", "":
		flagSet.Usage()
		os.Exit(2)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %q\n", cmd)
		flagSet.Usage()
		os.Exit(2)
	}
}

package main

import (
	"flag"
	"fmt"
)

// var myFlag = flag.Bool("flagname", false, "TODO...")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "usage: [command name] [TODO(fill this in)]")
		flag.PrintDefaults()
	}

	flag.Parse()

	// if flag.NArg() != 2 {
	// Uncomment to check number of args.
	// }

	// switch flag.Arg(0) {
	// Uncomment to switch on the first arg.
	// }
}

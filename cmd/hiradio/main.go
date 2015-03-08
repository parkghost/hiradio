// Command hiradio helps to play radio via Hichannel.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type command struct {
	name        string
	description string
	run         func(args []string)
}

var commands = []command{
	{"list", "List radio stations", listCmd},
	{"info", "Display radio information and program list", infoCmd},
	{"play", "Play radio on player", playCmd},
}

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `hiradio helps to play radio via Hichannel
Usage:

        hiradio [options] command [arg...]

The commands are:
`)
		for _, c := range commands {
			fmt.Fprintf(os.Stderr, "    %-24s %s\n", c.name, c.description)
		}
		fmt.Fprintln(os.Stderr, `
Use "hiradio command -h" for more information about a command.`)
		flag.PrintDefaults()
		os.Exit(1)
	}

}

func main() {
	log.SetFlags(0)
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
	}

	command := flag.Arg(0)
	for _, c := range commands {
		if c.name == command {
			c.run(flag.Args()[1:])
			return
		}
	}
	fmt.Fprintf(os.Stderr, "unknown command %q\n", command)
	fmt.Fprintln(os.Stderr, `Run "hiradio -h" for usage.`)
}

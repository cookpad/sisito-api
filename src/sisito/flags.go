package sisito

import (
	"flag"
	"fmt"
	"os"
)

var version string

const (
	DefaultConfig = "config.tml"
)

type Flags struct {
	Config string
}

func ParseFlag() (flags *Flags) {
	flags = &Flags{}
	var showVersion bool

	flag.StringVar(&flags.Config, "config", DefaultConfig, "config file path")
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.Parse()

	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	return
}

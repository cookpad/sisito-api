package sisito

import (
	"flag"
)

const (
	DefaultConfig = "config.tml"
)

type Flags struct {
	Config string
}

func ParseFlag() (flags *Flags) {
	flags = &Flags{}
	flag.StringVar(&flags.Config, "config", DefaultConfig, "config file path")
	flag.Parse()
	return
}

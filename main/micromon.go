package main

import (
	"flag"
	"github.com/Chostakovitch/micromon"
)

func main() {
	//Handle command-line flags
	testing := flag.Bool("test", false, "Set the flag to run tests")
	confPath := flag.String("c", "mm.conf", "Path to the configuration file")
	flag.Parse()

	//Run in test mode : assert tests
	if *testing {
		micromon.LaunchTests()
	} else {
		micromon.Start(*confPath)
	}
}

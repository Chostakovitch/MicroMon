package main

import (
	"fmt"
	"config"
	"log"
)

func main() {
	conf, err := config.FetchConfig("mm.conf")
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Printf("%+v", conf)
}

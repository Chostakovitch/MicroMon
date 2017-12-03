package main

import (
	"log"

	"config"
	"urlwatch"
)

func main() {
	conf, err := config.FetchConfig("mm.conf")
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf("%+v", conf)

	meta, err := urlwatch.Test("https://github.com")
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf("%v", meta)
}

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
	urlwatch.CheckUrl("https://www.datadoghq.com/")
	urlwatch.WatchWebsites(conf)
	for {

	}
}

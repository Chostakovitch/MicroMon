package main

import (
	"log"

	"config"
	"urlwatch"
	"reflect"
)

func main() {
	conf, err := config.FetchConfig("mm.conf")
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf("%+v", conf)

	//Get map of channels ; each channel receive data for a website
	chans := urlwatch.WatchWebsites(conf)

	//Build one SelectCase per channel, and forever listen to any data coming
	//Solution adapted from : https://stackoverflow.com/a/19992525
	cases := make([]reflect.SelectCase, len(chans))
	for i, ch := range chans {
		cases[i] = reflect.SelectCase{Chan: reflect.ValueOf(ch), Dir: reflect.SelectRecv}
	}
	for {
		_, value, _ := reflect.Select(cases)
		data := value.Interface().(urlwatch.MetaResponse)
		log.Print(data.Name)
	}
}

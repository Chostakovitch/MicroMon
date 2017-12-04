package urlwatch

import (
	"log"
	"time"

	"config"
)

//WatchWebsites takes a Config and check the URL at user-defined intervals.
//It returns a map associating the friendly name of the website and a channel which will receive MetaResponse each time a request is performed.
//Timers for intervals are launched in goroutines.
func WatchWebsites(conf config.Config) (map[string]chan MetaResponse) {
	chans := make(map[string]chan MetaResponse, len(conf.Websites))
	for name, website := range conf.Websites {
		//TODO changer ce magic number
		chans[name] = make(chan MetaResponse, 100)
		go func(name string, website config.Website, dataChan chan MetaResponse) {
			log.Printf("%v %v", name, website)
			for range time.Tick(time.Duration(website.Interval) * time.Second) {
				log.Printf("%v", name)
				feedChan(website.URL, dataChan)
			}
		}(name, website, chans[name])
	}
	return chans
}

//FeedChan takes an URL, check it, compute a MetaResponse and put it in a channel to make it compatible with the use of goroutines.
func feedChan(url string, data chan MetaResponse) {
	metaResp, err := CheckUrl(url)
	if err != nil {
		log.Fatalf("%v", err)
	}
	data <- metaResp
}

package urlwatch

import (
	"log"
	"time"

	"config"
)

//WatchWebsites takes a Config and check the URL at user-defined intervals.
//It returns a channel which will receive MetaResponse each time a request is performed.
//Timers for intervals are launched in goroutines.
func WatchWebsites(conf config.Config) (chan MetaResponse) {
	//TODO changer ce magic number
	ch := make(chan MetaResponse, 100)
	for name, website := range conf.Websites {
		go func(name string, website config.Website, dataChan chan MetaResponse) {
			for range time.Tick(time.Duration(website.Interval) * time.Second) {
				feedChan(website.URL, name, dataChan)
			}
		}(name, website, ch)
	}
	return ch
}

//FeedChan takes an url, check it, compute a MetaResponse and put it in a channel to make it compatible with the use of goroutines.
func feedChan(url string, name string, data chan MetaResponse) {
	metaResp, err := CheckUrl(url)
	if err != nil {
		log.Fatalf("%v", err)
	}
	metaResp.Name = name
	data <- metaResp
}

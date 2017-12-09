//urlwatch contains functions to watch URLs, i.e. check at regular intervals, using Config type as input.
//urlwatch defines the MetaResponse type, which holds websites responses' metadata.
package urlwatch

import (
	"log"
	"time"

	"config"
)

//WatchWebsites takes a Config and check the URLs at user-defined intervals.
//Returns a channel which will receive MetaResponse each time a request is performed.
//Each website has its own timed goroutine.
func WatchWebsites(conf config.Config) chan MetaResponse {
	//TODO changer ce magic number
	ch := make(chan MetaResponse, 100)
	for name, website := range conf.Websites {
		//Each website has an associated goroutine which perform a request
		//every X seconds, X user-defined, and send it to the global channel
		go func(name string, website config.Website, dataChan chan MetaResponse) {
			for range time.Tick(time.Duration(website.Interval) * time.Second) {
				feedChan(website.URL, name, dataChan)
			}
		}(name, website, ch)
	}
	return ch
}

//feedChan takes an url, check it, compute a MetaResponse and
//put it in a channel to make it compatible with the use of goroutines.
func feedChan(url string, name string, data chan MetaResponse) {
	metaResp, err := CheckUrl(url)
	if err != nil {
		log.Fatalf("%v", err)
	}
	//Add name to make MetaResponse independant
	metaResp.Name = name
	data <- metaResp
}

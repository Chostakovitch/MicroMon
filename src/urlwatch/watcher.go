//Package urlwatch contains functions to watch URLs, i.e. check them at regular intervals, using configuration as input.
//urlwatch defines the MetaResponse type, which holds websites responses' metadata.
package urlwatch

import (
	"log"
	"time"

	"config"
)

//WatchWebsites takes the app configuration and checks the URLs at user-defined intervals.
//It returns a channel which will receive MetaResponse each time a request is completed.
//Each website has its own timed goroutine.
func WatchWebsites(conf config.Config) chan MetaResponse {
	ch := make(chan MetaResponse, 100)
	for name, website := range conf.Websites {
		//Each website has an associated goroutine which perform a request
		//every X seconds, X user-defined, and send it to the global channel
		go func(name string, website config.Website, dataChan chan MetaResponse, timeout time.Duration) {
			for range time.Tick(time.Duration(website.Interval) * time.Second) {
				feedChan(website.URL, name, dataChan, timeout)
			}
		}(name, website, ch, time.Duration(conf.Timeout)*time.Second)
	}
	return ch
}

//feedChan takes an url, check it with a custom timeout, compute a MetaResponse and
//put it in a channel to make it compatible with the use of goroutines.
func feedChan(url string, name string, data chan MetaResponse, timeout time.Duration) {
	metaResp, err := CheckUrl(url, timeout)
	if err != nil {
		log.Fatalf("%v", err)
	}
	//Add name to make MetaResponse independent
	metaResp.Name = name
	data <- metaResp
}

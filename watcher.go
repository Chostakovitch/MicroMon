package micromon

import (
	"log"
	"net/http"
	"net/http/httptrace"
	"time"
)

//MetaResponse holds a website response's metadata, e.g. response code, response time, availibity, language...
type MetaResponse struct {
	URL              string
	Name             string
	Code             int
	wroteRequestTime time.Time
	RespDuration     time.Duration
	Timestamp        time.Time
	Available        bool
}

//WatchWebsites takes the app configuration and checks the URLs at user-defined intervals.
//It returns a channel which will receive MetaResponse each time a request is completed.
//Each website has its own timed goroutine.
func WatchWebsites(conf Config) chan MetaResponse {
	ch := make(chan MetaResponse, 100)
	for name, website := range conf.Websites {
		//Each website has an associated goroutine which perform a request
		//every X seconds, X user-defined, and send it to the global channel
		go func(name string, website Website, dataChan chan MetaResponse, timeout time.Duration) {
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

//CheckUrl produces a MetaResponse after visiting a given Website.
//So called "response time" is measured as the interval bewteen the start of server processing and the first byte received.
func CheckUrl(url string, timeout time.Duration) (MetaResponse, error) {
	meta := MetaResponse{URL: url}

	//New Client with low timeout
	client := http.Client{Timeout: timeout}

	//GET request with trace behaviour
	req, _ := http.NewRequest("GET", url, nil)
	req = withMetaResponse(req, &meta)

	//Perform GET request, feed MetaResponse and return values
	resp, err := client.Do(req)

	meta.Timestamp = time.Now()

	//We consider all errors as unavailability (if we only handle net.error Timeout error type, a non-existing URL throws an error)
	if err != nil {
		meta.Available = false
	} else {
		meta.Available = true
		meta.Code = resp.StatusCode
	}

	return meta, nil
}

//withMetaResponse adapts an HTTP Request to feed a MetaResponse object while performing request, thank to httptrace features.
//Returns a pointer to the augmented Request.
func withMetaResponse(req *http.Request, meta *MetaResponse) *http.Request {
	newReq := req.WithContext(
		httptrace.WithClientTrace(
			req.Context(),
			&httptrace.ClientTrace{
				//After DNS lookup and eventual TLS handshake : server starts processing
				WroteRequest: func(info httptrace.WroteRequestInfo) {
					meta.wroteRequestTime = time.Now()
				},
				//Server has processed and first byte is received : able to calculate accurate response-time
				GotFirstResponseByte: func() {
					meta.RespDuration = time.Now().Sub(meta.wroteRequestTime)
				},
			}),
	)
	return newReq
}

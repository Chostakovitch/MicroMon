//urlwatch contains functions to watch URLs, i.e. check at regular intervals, using Config type as input.
//urlwatch defines the MetaResponse type, which holds websites responses' metadata.
package urlwatch

import (
	"net/http"
	"net/http/httptrace"
	"time"

	"net"
)

//MetaResponse holds a website response's metadata, e.g. response code, response time, availibity, language...
type MetaResponse struct {
	URL string
	Name string
	Code      int
	wroteRequestTime time.Time
	RespDuration  time.Duration
	Timestamp time.Time
	Available bool
}

//checkUrl produces a MetaResponse after visiting a given Website.
//So called "response time" is measured as the interval bewteen the start of server processing and the first byte received.
func CheckUrl(url string) (MetaResponse, error) {
	meta := MetaResponse{URL: url}

	//New Client with low timeout
	//TODO put timeout in config
	client := http.Client{Timeout: 3 * time.Second}

	//GET request with trace behaviour
	req, _ := http.NewRequest("GET", url, nil)
	req = withMetaResponse(req, &meta)

	//Perform GET request, feed MetaResponse and return values
	resp, err := client.Do(req)

	meta.Timestamp = time.Now()
	//Timeout
	if err != nil {
		//Timeout error : we handle that one
		if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			meta.Available = false
		} else {
			return meta, err
		}
	} else {
		meta.Available = true
	}

	//HTTP code
	meta.Code = resp.StatusCode

	return meta, err
}

//withMetaResponse adapts an HTTP Request to feed a MetaResponse object while performing request, thank to httptrace features.
//Return a pointer to the augmented Request.
func withMetaResponse(req *http.Request, meta *MetaResponse) (*http.Request) {
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
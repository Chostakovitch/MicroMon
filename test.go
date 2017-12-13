package micromon

import (
	"net/http"
	"time"
)

//TestAlerting tests the alerting logic. It simulates a real word situation and configuration.
//It watches a local webserver for a while and compute availability. It also calls the alerting hook.
//At some point, it stops the webserver and repeat the same process.
//It checks if the hook returns a valid string which indicates that the web server is now marked as unavailable.
//It restart the web server and checks if the hook has marked the web server as recovered.
//If all theses conditions are met, it returns true ; false otherwise.
func TestAlerting() bool {
	//Create a local webserver
	http.HandleFunc("/", dummyResponse)
	srv := startHttpServer()

	//Create minimal configuration
	webserv := make(map[string]Website)
	webserv["localhost"] = Website{"http://localhost:8080", 1}
	conf := Config{Websites: webserv, Timeout: 3, AvailThreshold: 50}

	//Gather incoming MetaResponse
	ch := WatchWebsites(conf)
	data := NewSafeData()
	go func() {
		for {
			resp := <-ch
			data.Mux.Lock()
			data.Datas = append(data.Datas, resp)
			data.Mux.Unlock()
		}
	}()

	//Hook for alerting logic
	hook := AlertHook{}.GetHook(conf)

	//Local webserver is available, hook should not return anything
	time.Sleep(2 * time.Second)
	if getAvailStatus(data.Datas, conf, hook) != "" {
		return false
	}

	//We shutdown local webserver, hook should detect a new unavailability
	srv.Close()
	time.Sleep(2 * time.Second)
	if getAvailStatus(data.Datas, conf, hook) != "unavailable" {
		return false
	}

	//We restart local webserver, hook should detect recovery
	srv = startHttpServer()
	time.Sleep(2 * time.Second)
	defer srv.Close()
	if getAvailStatus(data.Datas, conf, hook) != "recovered" {
		return false
	}

	//Test is passed
	return true
}

//dummyResponse is a handler which responds with a 200 HTTP code for testing purpose
func dummyResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

//getAvailStatus computes availability metric and return alerting hook message
func getAvailStatus(resp []MetaResponse, conf Config, alertHook Hook) string {
	avail := Availability{}.Compute(resp)
	dummyWebs := make([]WebMetrics, 1)
	dummyWeb := WebMetrics{10, "localhost", make([]WebMetric, 1)}
	dummyWeb.Metrics = append(dummyWeb.Metrics, WebMetric{Availability{}, avail})
	dummyWebs = append(dummyWebs, dummyWeb)
	return alertHook(dummyWebs)
}

//startHttpServer starts a dummy Http server and returns a reference to it
func startHttpServer() *http.Server {
	srv := &http.Server{Addr: ":8080"}
	go func() {
		srv.ListenAndServe()
	}()
	return srv
}

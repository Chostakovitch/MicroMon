package micromon

import (
	"log"
	"time"
)

//Start... starts MicroMon with the configuration file which path is given in parameter.
func Start(path string) {
	//Get configuration from file
	conf := GetConfig(path)

	//Get channel for receiving website response data
	ch := WatchWebsites(conf)

	//Create data structure for holding response data
	datas := NewRespMap(len(conf.Websites))
	for k, _ := range conf.Websites {
		datas[k] = NewSafeData()
	}

	//Configure metrics, hook and reporter
	metrics := GetMetrics(conf)
	reporter := GetReporter(conf)
	hooks := GetHooks(conf)

	//Compute and write metrics every 10 seconds
	go func() {
		i := 0
		for range time.Tick(10 * time.Second) {
			i++
			res := make([][]WebMetrics, 0)
			//Metrics for the last 2 minutes and the last 10 minutes
			res = append(res, (&datas).ComputeMetrics(metrics, 2))

			//We apply hooks only once (avoiding repeating logging)... TODO May be improved!
			applyHooks(res, hooks)
			res = append(res, (&datas).ComputeMetrics(metrics, 10))

			//Metrics for the last hour are reported every minute
			if i%6 == 0 {
				res = append(res, (&datas).ComputeMetrics(metrics, 60))
			}

			reportResults(res, reporter)
		}
	}()

	//Forever listen to data coming from channel : sequential access to datas variable
	for {
		data := <-ch
		name := data.Name
		datas[name].Mux.Lock()
		datas[name].Datas = append(datas[name].Datas, data)
		datas[name].Mux.Unlock()
	}
}

//launchTests performs tests against the application logic and report results.
func LaunchTests() {
	log.Printf("Starting tests...")
	if TestAlerting() {
		log.Print("Alerting test successfully passed !")
	} else {
		log.Fatalf("Alert test failed !")
	}
	log.Printf("All tests passed !")
}

//reportResult takes a slice of []WebMetrics and report information for each inner slice.
func reportResults(metrics [][]WebMetrics, reporter Reporter) {
	for _, v := range metrics {
		reporter.Report(v)
	}
}

//applyHook takes a slice of []WebMetrics and call hooks for each inner slice.
func applyHooks(metrics [][]WebMetrics, hooks []Hook) {
	for _, v := range metrics {
		for _, h := range hooks {
			h(v)
		}
	}
}

//main contains methods to orchestrate the application logic.
//First, it uses the config package to fetch the user configuration.
//Second, it uses the urlwatch package to get a channel of HTTP responses.
//Third, it uses the metric package to store HTTP responses and launch/collect a bunch of metrics computation.
//Fourth, it uses the hook package to launch intermediate work on metrics.
//Finally, it uses the report package to format and write metrics.
//Optionally, it uses the test package to launch some tests against the app logic.
package main

import (
	"flag"
	"hook"
	"log"
	"test"
	"time"

	"metric"
	"report"
	"urlwatch"
)

func main() {
	//Handle command-line flags
	testing := flag.Bool("test", false, "Set the flag to run tests")
	confPath := flag.String("c", "mm.conf", "Path to the configuration file")
	flag.Parse()

	//Run in test mode : assert tests
	if *testing {
		launchTests()
	} else {
		start(*confPath)
	}
}

//start... starts MicroMon with the configuration file which path is given in parameter.
func start(path string) {
	//Get configuration from file
	conf := getConfig(path)

	//Get channel for receiving website response data
	ch := urlwatch.WatchWebsites(conf)

	//Create data structure for holding response data
	datas := metric.NewRespMap(len(conf.Websites))
	for k, _ := range conf.Websites {
		datas[k] = metric.NewSafeData()
	}

	//Configure metrics, hook and reporter
	metrics := getMetrics(conf)
	reporter := getReporter(conf)
	hooks := getHooks(conf)

	//Compute and write metrics every 10 seconds / 1 minute
	go func() {
		i := 0
		for range time.Tick(10 * time.Second) {
			i++
			res := make([][]metric.WebMetrics, 0)
			//Metrics for the last 2 minutes and the last 10 minutes
			res = append(res, (&datas).ComputeMetrics(metrics, 2))

			//We apply hooks only once (avoiding repeating logging)... TODO May be improved!
			applyHooks(res, hooks)
			res = append(res, (&datas).ComputeMetrics(metrics, 10))

			//Metrics for the last hour every minute
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

//launchTests performs tests against the application logic and report result.
func launchTests() {
	log.Printf("Starting tests...")
	if test.TestAlerting() {
		log.Print("Alerting test successfully passed !")
	} else {
		log.Fatalf("Alert test failed !")
	}
	log.Printf("All tests passed !")
}

//reportResult takes a slice of []WebMetrics and report each inner slice.
func reportResults(metrics [][]metric.WebMetrics, reporter report.Reporter) {
	for _, v := range metrics {
		reporter.Report(v)
	}
}

//applyHook takes a slice of []WebMetrics and call hooks for each inner slice.
func applyHooks(metrics [][]metric.WebMetrics, hooks []hook.Hook) {
	for _, v := range metrics {
		for _, h := range hooks {
			h(v)
		}
	}
}

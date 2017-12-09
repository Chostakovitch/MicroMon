package main

import (
	"log"
	"time"

	"config"
	"metric"
	"report"
	"urlwatch"
)

func getConfig() config.Config {
	conf, err := config.FetchConfig("mm.conf")
	if err != nil {
		log.Fatalf("%v", err)
	}
	return conf
}

func main() {
	//Get configuration from file
	conf := getConfig()

	//Get channel for website response data
	ch := urlwatch.WatchWebsites(conf)

	//Create data structure for response data
	datas := metric.NewRespMap(len(conf.Websites))
	for k, _ := range conf.Websites {
		datas[k] = metric.NewSafeData()
	}

	//Configure metrics and reporter
	metrics := []metric.Metric{metric.AvgRespTime{}, metric.MaxRespTime{}, metric.CodeCount{}, metric.Availability{}}
	reporter := report.NewReporter(report.DefaultLogger(), report.DefaultFormatter{})

	//Compute and write metrics every 10 seconds / 1 minute
	go func() {
		i := 0
		for range time.Tick(10 * time.Second) {
			i++
			reporter.Report(metrics, (&datas).ComputeMetrics(metrics, 2), 2)
			reporter.Report(metrics, (&datas).ComputeMetrics(metrics, 2), 10)
			if i%6 == 0 {
				reporter.Report(metrics, (&datas).ComputeMetrics(metrics, 2), 60)
			}
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

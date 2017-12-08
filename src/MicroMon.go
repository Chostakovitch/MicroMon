package main

import (
	"log"
	"time"

	"config"
	"report"
	"urlwatch"
	"metric"
)

func main() {
	//TODO splitter main en plusieurs fonctions
	conf, err := config.FetchConfig("mm.conf")
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf("%+v", conf)

	//Get channel for website response data
	ch := urlwatch.WatchWebsites(conf)

	datas := metric.NewRespMap(len(conf.Websites))
	for k, _ := range conf.Websites {
		datas[k] = metric.NewSafeData()
	}


	metrics := []metric.Metric{metric.AvgRespTime{}, metric.MaxRespTime{}, metric.CodeCount{}, metric.Availibility{}}
	reporter := report.NewReporter(report.DefaultLogger(), report.DefaultFormatter{})


	go func() {
		for range time.Tick(10 * time.Second) {
			res := (&datas).ComputeMetrics(metrics, 2)
			for k, v := range res {
				reporter.Report(metrics, v, k, 2)
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

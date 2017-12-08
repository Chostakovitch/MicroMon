package main

import (
	"log"

	"config"
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

	go metric.CalculateMetrics(&datas, []metric.Metric{metric.AvgRespTime{}})

	//Forever listen to data coming from channel : sequential access to datas variable
	for {
		data := <-ch
		name := data.Name
		datas[name].Mux.Lock()
		datas[name].Datas = append(datas[name].Datas, data)
		datas[name].Mux.Unlock()
	}
}

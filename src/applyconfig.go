package main

import (
	"config"
	"hook"
	"log"
	"metric"
	"report"
)

//getConfig returns a Config fetched from the path given in parameter
func getConfig(path string) config.Config {
	conf, err := config.FetchConfig(path)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return conf
}

//getMetrics returns instances of Metric from the configuration
func getMetrics(conf config.Config) []metric.Metric {
	metrics := make([]metric.Metric, 0)
	//Try to instantiate each metric
	for _, v := range conf.Metrics {
		met, err := metric.GetMetric(v)
		if err != nil {
			log.Printf("Warning : %v", err)
		} else {
			metrics = append(metrics, met)
		}
	}
	return metrics
}

//getMetrics returns instances of Hook from the configuration
func getHooks(conf config.Config) []hook.Hook {
	hooks := make([]hook.Hook, 0)
	//Try to instantiate each hook : get the Hooker and the Hook closure with closed-config.
	for _, v := range conf.Hooks {
		h, err := hook.GetHook(v, conf)
		if err != nil {
			log.Printf("Warning : %v", err)
		} else {
			hooks = append(hooks, h)
		}
	}
	return hooks
}

//getReporter builds a reporter from the configuration
func getReporter(conf config.Config) report.Reporter {
	var f report.Formatter
	var l *log.Logger

	//Get formatter
	switch conf.Format {
	case "xml":
		f = report.XMLFormatter{}
	default:
		f = report.DefaultFormatter{}
	}

	//Get logger
	switch conf.Output {
	case "":
		l = report.DefaultLogger()
	default:
		l = report.FileLogger(conf.Output)
	}

	return report.NewReporter(l, f)
}

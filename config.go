package micromon

import (
	"io/ioutil"
	"log"

	"yaml.v2"
)

//Config is a struct which mirrors the structure of the YAML configuration file.
//It contains user customisable parameters, such as websites to visit and metrics to compute.
type Config struct {
	Websites        map[string]Website
	DefaultInterval int
	Timeout         int
	AvailThreshold  int
	Metrics         []string
	Hooks           []string
	Format          string
	Output          string
}

//Website wraps an URL and a check interval.
type Website struct {
	URL      string
	Interval int
}

//FetchConfig parses a YAML file which reflects MicroMon's configuration.
//It takes an input path and return a Config object - or an error.
func FetchConfig(path string) (Config, error) {
	conf := Config{}

	//Read and decode configuration from file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return conf, err
	}
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return conf, err
	}

	//Set default interval for unspecified check intervals
	for k, v := range conf.Websites {
		//Workaround because we cannot assign to struct field in map, so copy struct, make change, assign strut
		tmp := conf.Websites[k]
		if v.Interval == 0 {
			tmp.Interval = conf.DefaultInterval
			conf.Websites[k] = tmp
		}
	}
	return conf, err
}

//GetConfig returns a Config fetched from the path given in parameter.
func GetConfig(path string) Config {
	conf, err := FetchConfig(path)
	if err != nil {
		log.Fatalf("Error fetching configuration : %v", err)
	}
	return conf
}

//GetMetrics returns instances of Metric from the configuration.
func GetMetrics(conf Config) []Metric {
	metrics := make([]Metric, 0)
	//Try to instantiate each metric
	for _, v := range conf.Metrics {
		met, err := GetMetric(v)
		if err != nil {
			log.Printf("Warning : %v", err)
		} else {
			metrics = append(metrics, met)
		}
	}
	return metrics
}

//GetHooks returns instances of Hook from the configuration.
func GetHooks(conf Config) []Hook {
	hooks := make([]Hook, 0)
	//Try to instantiate each hook : get the Hooker and the Hook closure with closed-config.
	for _, v := range conf.Hooks {
		h, err := GetHook(v, conf)
		if err != nil {
			log.Printf("Warning : %v", err)
		} else {
			hooks = append(hooks, h)
		}
	}
	return hooks
}

//GetReporter builds a reporter from the configuration
func GetReporter(conf Config) Reporter {
	var f Formatter
	var l *log.Logger

	//Get formatter
	switch conf.Format {
	case "xml":
		f = XMLFormatter{}
	default:
		f = DefaultFormatter{}
	}

	//Get logger
	switch conf.Output {
	case "":
		l = DefaultLogger()
	default:
		l = FileLogger(conf.Output)
	}

	return NewReporter(l, f)
}

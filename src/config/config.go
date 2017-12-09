//config contains utilities to retrieve and build MicroMon configuration.
package config

import (
	"io/ioutil"
	"yaml.v2"
)

//Config is just a set of websites to check and a default interval.
type Config struct {
	Websites        map[string]Website
	DefaultInterval int
}

//Website is an URL and a check interval.
type Website struct {
	URL      string
	Interval int
}

//FetchConfig parses a YAML file which defines websites to visit and check invervals.
//Takes an input path and return a Config object - or an error.
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

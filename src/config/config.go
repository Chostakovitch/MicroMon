//config contains utilities to retrieve and build MicroMon configuration.
package config

import (
	"yaml.v2"
	"io/ioutil"
)

//A Config is just a set of websites to check and a default interval.
type Config struct {
	Websites        map[string]Website
	DefaultInterval int
}

//A Website is an URL and a check interval.
type Website struct {
	URL      string
	Interval int
}


//FetchConfig parses a YAML file which defines the websites to visit and the check inverval.
//Takes an input path and return a Config object - or an error.
func FetchConfig(path string) (Config, error) {
	conf := Config{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return conf, err
	}

	err = yaml.Unmarshal(data, &conf)
	return conf, err
}

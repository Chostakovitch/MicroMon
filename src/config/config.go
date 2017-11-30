package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Websites        map[string]Website
	DefaultInterval int
}

type Website struct {
	URL      string
	Interval int
}

func FetchConfig(path string) (Config, error) {
	conf := Config{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return conf, err
	}

	err = yaml.Unmarshal(data, &conf)
	return conf, err
}

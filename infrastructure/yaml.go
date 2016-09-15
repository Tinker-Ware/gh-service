package infrastructure

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Configuration stores the fields to configure the application
type Configuration struct {
	Port         string   `yaml:"port"`
	ClientID     string   `yaml:"clientID"`
	ClientSecret string   `yaml:"clientSecret"`
	Salt         string   `yaml:"salt"`
	Scopes       []string `yaml:"scopes,flow"`
	APIHost      string
}

// GetConfiguration returns the configuration stored in a file
func GetConfiguration(path string) (*Configuration, error) {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	conf := &Configuration{}

	err = yaml.Unmarshal(data, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil

}

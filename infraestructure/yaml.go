package infraestructure

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Port         string   `yaml:"port"`
	ClientID     string   `yaml:"clientID"`
	ClientSecret string   `yaml:"clientSecret"`
	Scopes       []string `yaml:"scopes,flow"`
}

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

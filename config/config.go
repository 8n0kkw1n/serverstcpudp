package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config - struct config
type Config struct {
	TCP  TCPConfig  `yaml:"tcp"`
	UDP  UDPConfig  `yaml:"udp"`
	HTTP HTTPConfig `yaml:"http"`
}

// TCPConfig - structure tcp connect
type TCPConfig struct {
	Host string `yaml:"host"`
	Port uint   `yaml:"port"`
}

// UDPConfig - structure udp connect
type UDPConfig struct {
	Host string `yaml:"host"`
	Port uint16 `yaml:"port"`
}

// HTTPConfig - structure HTTP connect
type HTTPConfig struct {
	Host string `yaml:"host"`
	Port uint16 `yaml:"port"`
}

// NewConfig - new config file create
func NewConfig(path string) (Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	conf := Config{}
	if err := yaml.Unmarshal(file, &conf); err != nil {
		return Config{}, err
	}
	return conf, nil
}

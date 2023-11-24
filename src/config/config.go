package config

import (
	"fmt"
	"io/ioutil"
)

func ReadConfigFromFile(filepath string) (*Config, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %v", err)
	}
	parser := NewMultiParser(&YAMLParser{}, &JSONParser{})
	config := &Config{}
	err = parser.Parse(file, config)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML file: %v", err)
	}
	err = config.Validate()
	if err != nil {
		return nil, fmt.Errorf("error validating YAML file: %v", err)
	}
	return config, nil
}

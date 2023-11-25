package config

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func LoadConfig() (*Config, error) {
	path, ok := os.LookupEnv("VELOCITY_CONFIG")
	if !ok {
		path = "velocity.yml"
	}
	// Switch statement to find out what filepath starts with
	// if filepath starts with http:// or https://
	// then use ReadConfigFromURL
	// else use ReadConfigFromFile
	switch {
	case len(path) > 8 && (path[:7] == "http://" || path[:8] == "https://"):
		return readConfigFromURL(path)
	default:
		return readConfigFromFile(path)
	}
}

func readConfigFromURL(url string) (*Config, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error retrieving file '%s': status code '%d' error '%v'", url, response.StatusCode, err)
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return parseConfig(bytes, NewMultiParser(&YAMLParser{}, &JSONParser{}))

}

func readConfigFromFile(filepath string) (*Config, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("error reading file '%s': %v", filepath, err)
	}
	return parseConfig(file, NewMultiParser(&YAMLParser{}, &JSONParser{}))
}

func parseConfig(config []byte, parser MultiParser) (*Config, error) {
	c := &Config{}
	err := parser.Parse(config, c)
	if err != nil {
		return nil, fmt.Errorf("error parsing config: %v", err)
	}
	err = c.Populate()
	if err != nil {
		return nil, fmt.Errorf("error populating config: %v", err)
	}
	err = c.Validate()
	if err != nil {
		return nil, fmt.Errorf("error validating config: %v", err)
	}
	return c, nil
}

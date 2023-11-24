package config

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func (c *Config) GetWorkflow(w string) (YAMLWorkflow, error) {
	for title, workflow := range c.Workflows {
		if title == w {
			return workflow, nil
		}
	}
	return YAMLWorkflow{}, fmt.Errorf("workflow %s not found", w)
}

func (c *Config) GetWorkflowNames() []string {
	var names []string
	for name := range c.Workflows {
		names = append(names, name)
	}
	return names
}

func LoadConfig() (*Config, error) {
	filePath, ok := os.LookupEnv("VELOCITY_CONFIG_FILE")
	if !ok {
		filePath = "velocity.yml"
	}
	// Switch statement to find out what filepath starts with
	// if filepath starts with http:// or https://
	// then use ReadConfigFromURL
	// else use ReadConfigFromFile
	switch {
	case filePath[:7] == "http://" || filePath[:8] == "https://":
		return ReadConfigFromURL(filePath)
	default:
		return ReadConfigFromFile(filePath)
	}
}

func ReadConfigFromURL(url string) (*Config, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to download file. Status code: %d. %v", response.StatusCode, err)
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return ParseConfig(bytes, NewMultiParser(&YAMLParser{}, &JSONParser{}))

}

func ReadConfigFromFile(filepath string) (*Config, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %v", err)
	}
	return ParseConfig(file, NewMultiParser(&YAMLParser{}, &JSONParser{}))
}

func ParseConfig(config []byte, parser MultiParser) (*Config, error) {
	c := &Config{}
	err := parser.Parse(config, c)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML file: %v", err)
	}
	err = c.Validate()
	if err != nil {
		return nil, fmt.Errorf("error validating YAML file: %v", err)
	}
	err = c.Populate()
	if err != nil {
		return nil, fmt.Errorf("error populating YAML file: %v", err)
	}
	return c, nil
}

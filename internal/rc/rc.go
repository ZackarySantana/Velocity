package rc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"syscall"

	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

type RuntimeConfigurations struct {
	Servers []RuntimeConfiguration `yaml:"servers"`
}

type RuntimeConfiguration struct {
	Server string `yaml:"server"`
	APIKey string `yaml:"api_key"`
}

func GetRuntimeConfigurations() (*RuntimeConfigurations, error) {
	configFile, err := getRCFile()
	if err != nil {
		return nil, err
	}

	return readRCFromFile(configFile)
}

func GetRuntimeConfiguration(server string) (*RuntimeConfiguration, error) {
	configFile, err := getRCFile()
	if err != nil {
		return nil, err
	}

	rcs, err := readRCFromFile(configFile)
	if err != nil {
		return nil, err
	}
	if rcs == nil {
		return nil, errors.New("no runtime configurations found")
	}

	for _, r := range rcs.Servers {
		if r.Server == server {
			return &r, nil
		}
	}

	return nil, fmt.Errorf("server %s not found", server)
}

func SetRuntimeConfigurations(config *RuntimeConfigurations) error {
	configFile, err := getRCFile()
	if err != nil {
		return err
	}

	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configFile, yamlData, 0644)
}

func AskForRuntimeConfiguration(server string) (*RuntimeConfiguration, error) {
	rcs, err := GetRuntimeConfigurations()
	if err != nil || rcs == nil {
		rcs = &RuntimeConfigurations{}
	}

	rc := RuntimeConfiguration{
		Server: server,
	}
	fmt.Print("Enter API key: ")
	rc.APIKey, err = promptSecret()
	if err != nil {
		return nil, err
	}

	rcs.Servers = append(rcs.Servers, rc)
	// Write updated configurations to file
	err = SetRuntimeConfigurations(rcs)
	if err != nil {
		return nil, err
	}

	fmt.Println("API key and server added to the configuration file.")
	return &rc, nil
}

func promptSecret() (string, error) {
	pw, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	fmt.Println()
	return string(pw), nil
}

func readRCFromFile(filePath string) (*RuntimeConfigurations, error) {
	var config *RuntimeConfigurations

	yamlData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return config, yaml.Unmarshal(yamlData, &config)
}

func getRCFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return filepath.Join(usr.HomeDir, ".velocityrc"), nil
}

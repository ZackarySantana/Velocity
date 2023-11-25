package config

import (
	"encoding/json"
	"errors"

	"gopkg.in/yaml.v3"
)

type MultiParser []Parser
type Parser interface {
	// Parse parses the given byte slice and stores the result in the value pointed to by out. Out should be a pointer to a struct.
	Parse(in []byte, out interface{}) error
}

type YAMLParser struct{}

func (p *YAMLParser) Parse(data []byte, out interface{}) error {
	return yaml.Unmarshal(data, out)
}

type JSONParser struct{}

func (p *JSONParser) Parse(data []byte, out interface{}) error {
	return json.Unmarshal(data, out)
}

func NewMultiParser(parsers ...Parser) MultiParser {
	return parsers
}

func (p MultiParser) Parse(data []byte, out interface{}) error {
	for _, parser := range p {
		err := parser.Parse(data, out)
		if err == nil {
			return nil
		}
	}
	return errors.New("unable to parse data")
}

// Code generated from Pkl module `test_construct`. DO NOT EDIT.
package test

import "github.com/zackarysantana/velocity/gen/pkl/primitives/command"

type Test struct {
	Name string `pkl:"name" bson:"name,omitempty" json:"name,omitempty" yaml:"name,omitempty"`

	Commands []command.Command `pkl:"commands" bson:"commands,omitempty" json:"commands,omitempty" yaml:"commands,omitempty"`
}

package entities

import (
	"github.com/zackarysantana/velocity/src/entities/image"
	"github.com/zackarysantana/velocity/src/entities/job"
	"github.com/zackarysantana/velocity/src/entities/routine"
	"github.com/zackarysantana/velocity/src/entities/test"
)

// ConfigEntity is not meant to be stored in the database
// but as a DTO for an entity representation of the configuration
// file.
// This will be constructed images/tests -> jobs -> routines.
type ConfigEntity struct {
	Images   []*image.Image
	Tests    []*test.Test
	Jobs     []*job.Job
	Routines []*routine.Routine
}

func (c *ConfigEntity) Merge(other *ConfigEntity) {
	c.Images = append(c.Images, other.Images...)
	c.Tests = append(c.Tests, other.Tests...)
	c.Jobs = append(c.Jobs, other.Jobs...)
	c.Routines = append(c.Routines, other.Routines...)
}

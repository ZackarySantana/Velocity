// Code generated from Pkl module `velocity`. DO NOT EDIT.
package velocity

import (
	"context"

	"github.com/apple/pkl-go/pkl"
	"github.com/zackarysantana/velocity/gen/pkl/constructs/build"
	"github.com/zackarysantana/velocity/gen/pkl/constructs/deployment"
	"github.com/zackarysantana/velocity/gen/pkl/constructs/test"
	"github.com/zackarysantana/velocity/gen/pkl/constructs/workflow"
	"github.com/zackarysantana/velocity/gen/pkl/primitives/runtime"
)

type Velocity struct {
	Tests []*test.Test `pkl:"tests" bson:"tests,omitempty" json:"tests,omitempty" yaml:"tests,omitempty"`

	Runtimes []runtime.Runtime `pkl:"runtimes" bson:"runtimes,omitempty" json:"runtimes,omitempty" yaml:"runtimes,omitempty"`

	Workflows []*workflow.Workflow `pkl:"workflows" bson:"workflows,omitempty" json:"workflows,omitempty" yaml:"workflows,omitempty"`

	Builds []*build.Build `pkl:"builds" bson:"builds,omitempty" json:"builds,omitempty" yaml:"builds,omitempty"`

	Deployments []*deployment.Deployment `pkl:"deployments" bson:"deployments,omitempty" json:"deployments,omitempty" yaml:"deployments,omitempty"`

	Output any `pkl:"output" bson:"output,omitempty" json:"output,omitempty" yaml:"output,omitempty"`
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Velocity
func LoadFromPath(ctx context.Context, path string) (ret *Velocity, err error) {
	evaluator, err := pkl.NewEvaluator(ctx, pkl.PreconfiguredOptions)
	if err != nil {
		return nil, err
	}
	defer func() {
		cerr := evaluator.Close()
		if err == nil {
			err = cerr
		}
	}()
	ret, err = Load(ctx, evaluator, pkl.FileSource(path))
	return ret, err
}

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Velocity
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (*Velocity, error) {
	var ret Velocity
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

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
	Tests []*test.Test `pkl:"tests"`

	Runtimes []runtime.Runtime `pkl:"runtimes"`

	Workflows []*workflow.Workflow `pkl:"workflows"`

	Builds []*build.Build `pkl:"builds"`

	Deployments []*deployment.Deployment `pkl:"deployments"`

	Output any `pkl:"output"`
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

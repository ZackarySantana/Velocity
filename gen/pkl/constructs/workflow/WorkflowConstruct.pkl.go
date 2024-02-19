// Code generated from Pkl module `workflow_construct`. DO NOT EDIT.
package workflow

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type WorkflowConstruct struct {
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a WorkflowConstruct
func LoadFromPath(ctx context.Context, path string) (ret *WorkflowConstruct, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a WorkflowConstruct
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (*WorkflowConstruct, error) {
	var ret WorkflowConstruct
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

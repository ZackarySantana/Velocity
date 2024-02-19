// Code generated from Pkl module `request_prebuilt`. DO NOT EDIT.
package request

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type RequestPrebuilt struct {
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a RequestPrebuilt
func LoadFromPath(ctx context.Context, path string) (ret *RequestPrebuilt, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a RequestPrebuilt
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (*RequestPrebuilt, error) {
	var ret RequestPrebuilt
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

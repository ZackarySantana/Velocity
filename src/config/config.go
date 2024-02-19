package config

import (
	"context"
	"errors"
	"strings"

	"github.com/apple/pkl-go/pkl"
	"github.com/zackarysantana/velocity/gen/pkl/velocity"
)

func Load(ctx context.Context, path string) (*velocity.Velocity, error) {
	if !strings.HasSuffix(path, ".pkl") {
		return nil, errors.New("unsupported file type for configuration file (.pkl only)")
	}
	var source *pkl.ModuleSource
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		source = pkl.UriSource(path)
	} else {
		source = pkl.FileSource(path)
	}

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

	// Check validation through output.value.
	_, err = evaluator.EvaluateExpressionRaw(ctx, source, "output.value")
	if err != nil {
		return nil, err
	}
	// Evaluate the module and return the result.
	var v velocity.Velocity
	if err := evaluator.EvaluateModule(ctx, source, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

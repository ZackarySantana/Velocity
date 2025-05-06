package velocity

import (
	"context"
	"fmt"
)

type linkClient struct {
	base string
}

func (c *linkClient) Routine(ctx context.Context, id string) string {
	return fmt.Sprintf("%s/routine/%v", c.base, id)
}

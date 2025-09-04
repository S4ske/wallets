package closer

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

type Closer struct {
	mu        sync.Mutex
	functions []func(context.Context) error
}

func (c *Closer) Add(f func(context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.functions = append(c.functions, f)
}

func (c *Closer) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var (
		msgs     = make([]string, 0, len(c.functions))
		complete = make(chan struct{}, 1)
	)

	go func() {
		for _, f := range c.functions {
			if err := f(ctx); err != nil {
				msgs = append(msgs, fmt.Sprintf("%v", err))
			}
		}

		complete <- struct{}{}
	}()

	select {
	case <-complete:
		break
	case <-ctx.Done():
		return fmt.Errorf("shutdown cancelled: %w", ctx.Err())
	}

	if len(msgs) > 0 {
		return fmt.Errorf(
			"shutdown finished with error(s): \n%s",
			strings.Join(msgs, "\n"),
		)
	}

	return nil
}

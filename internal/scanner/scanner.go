package scanner

import (
	"context"
	"fmt"
	"sync"

	"github.com/leylmordor/secretwatch/internal/store"
)

type Result struct {
	Secrets []*store.Secret
	Errors  []error
}

func Run(ctx context.Context, stores []store.Store) Result {
	var (
		mu      sync.Mutex
		result  Result
		wg      sync.WaitGroup
	)

	for _, s := range stores {
		wg.Add(1)
		go func(s store.Store) {
			defer wg.Done()
			secrets, err := s.List(ctx)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("[%s] %w", s.Type(), err))
				return
			}
			result.Secrets = append(result.Secrets, secrets...)
		}(s)
	}

	wg.Wait()
	return result
}

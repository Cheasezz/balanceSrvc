package tests

import "sync"

func runConcurrent(n int, fn func() error) []error {
	var wg sync.WaitGroup
	var start sync.WaitGroup

	errors := make([]error, 0, n)
	var mu sync.Mutex

	wg.Add(n)
	start.Add(1)

	for range n {
		go func() {
			defer wg.Done()

			start.Wait()

			if err := fn(); err != nil {
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
			}
		}()
	}

	start.Done()
	wg.Wait()

	return errors
}

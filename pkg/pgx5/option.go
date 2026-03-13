package pgx5

import "time"

type Option func(*Pgx)

func MinPoolSize(size int) Option {
	return func(c *Pgx) {
		c.minPoolSize = size
	}
}

func MaxPoolSize(size int) Option {
	return func(c *Pgx) {
		c.maxPoolSize = size
	}
}

func ConnAttempts(attempts int) Option {
	return func(c *Pgx) {
		c.connAttempts = attempts
	}
}

func connAttemptsDelay(timeout time.Duration) Option {
	return func(c *Pgx) {
		c.connAttemptsDelay = timeout
	}
}

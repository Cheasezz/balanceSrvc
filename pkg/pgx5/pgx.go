package pgx5

import (
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	ErrParseCfg      = "parse config error"
	ErrNewWithConfig = "newWithConfig error"
	ErrPing          = "ping db error"
)

const (
	_defaultMinPoolSize       = 5
	_defaultMaxPoolSize       = 10
	_defaultConnAttempts      = 10
	_defaultConnAttemptsDelay = time.Second
)

type Pgx struct {
	minPoolSize       int
	maxPoolSize       int
	connAttempts      int
	connAttemptsDelay time.Duration

	Scany *pgxscan.API
	Pool  *pgxpool.Pool
}

func New(url string, opts ...Option) (*Pgx, error) {
	const op = "pgx5.New"

	pg := &Pgx{
		minPoolSize:       _defaultMaxPoolSize,
		maxPoolSize:       _defaultMaxPoolSize,
		connAttempts:      _defaultConnAttempts,
		connAttemptsDelay: _defaultConnAttemptsDelay,
		Scany:             pgxscan.DefaultAPI,
	}

	for _, opt := range opts {
		opt(pg)
	}

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("op=%s: %s: %w", op, ErrParseCfg, err)
	}

	poolConfig.MinConns = int32(pg.minPoolSize)
	poolConfig.MaxConns = int32(pg.maxPoolSize)

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		fmt.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connAttemptsDelay)

		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("op=%s: %s; %w", op, ErrNewWithConfig, err)
	}

	if err = pg.Pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("op=%s: %s: %w", op, ErrPing, err)
	}

	fmt.Printf("Postgres connected, %d connAttempts remaining\n", pg.connAttempts)

	return pg, nil
}

func (p *Pgx) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

package repo

import (
	"context"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	systemrepo "github.com/Cheasezz/balanceSrvc/internal/repo/system"
)

type System interface {
	Transaction(c context.Context, req *core.SystemTransaction) error
}

type Repo struct {
	System System
}

func New() *Repo {
	return &Repo{
		System: systemrepo.New(),
	}
}

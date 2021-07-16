package memory

import (
	"context"

	"github.com/Tra-Dew/users/pkg/core"
	"github.com/Tra-Dew/users/pkg/users"
)

type repositoryInMemory struct {
	data map[string]*users.User
}

// NewRepository ...
func NewRepository() users.Repository {
	return &repositoryInMemory{
		data: make(map[string]*users.User),
	}
}

// Insert ...
func (r *repositoryInMemory) Insert(ctx context.Context, u *users.User) error {
	r.data[u.ID] = u
	return nil
}

// GetByEmail ...
func (r *repositoryInMemory) GetByEmail(ctx context.Context, id string) (*users.User, error) {
	u, exists := r.data[id]
	if !exists {
		return nil, core.ErrNotFound
	}

	return u, nil
}

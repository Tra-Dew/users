package memory

import (
	"context"

	"github.com/d-leme/tradew-users/pkg/core"
	"github.com/d-leme/tradew-users/pkg/users"
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
func (r *repositoryInMemory) GetByEmail(ctx context.Context, email string) (*users.User, error) {

	for _, user := range r.data {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, core.ErrNotFound
}

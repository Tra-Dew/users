package mock

import (
	"context"

	"github.com/Tra-Dew/users/pkg/users"
	"github.com/stretchr/testify/mock"
)

// RepositoryMock ...
type RepositoryMock struct {
	mock.Mock
}

// NewRepository ...
func NewRepository() users.Repository {
	return &RepositoryMock{}
}

// Insert ...
func (r *RepositoryMock) Insert(ctx context.Context, p *users.User) error {
	args := r.Mock.Called()

	arg0 := args.Get(0)
	if arg0 != nil {
		return arg0.(error)
	}

	return nil
}

// GetByEmail ...
func (r *RepositoryMock) GetByEmail(ctx context.Context, email string) (*users.User, error) {
	args := r.Mock.Called(email)

	arg0 := args.Get(0)
	if arg0 != nil {
		return arg0.(*users.User), nil
	}

	return nil, args.Get(1).(error)
}

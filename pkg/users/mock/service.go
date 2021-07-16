package mock

import (
	"context"

	"github.com/Tra-Dew/users/pkg/users"
	"github.com/stretchr/testify/mock"
)

// ServiceMock ...
type ServiceMock struct {
	mock.Mock
}

// NewService ...
func NewService() users.Service {
	return &ServiceMock{}
}

// Create ...
func (r *ServiceMock) Create(ctx context.Context, correlationID string, req *users.CreateUserRequest) (*users.CreateUserResponse, error) {
	args := r.Mock.Called()

	arg0 := args.Get(0)
	if arg0 != nil {
		return arg0.(*users.CreateUserResponse), nil
	}

	arg1 := args.Get(1)

	return nil, arg1.(error)
}

// Login ...
func (r *ServiceMock) Login(ctx context.Context, correlationID string, req *users.LoginRequest) (*users.LoginResponse, error) {
	args := r.Mock.Called()

	arg0 := args.Get(0)
	if arg0 != nil {
		return arg0.(*users.LoginResponse), nil
	}

	arg1 := args.Get(1)

	return nil, arg1.(error)
}

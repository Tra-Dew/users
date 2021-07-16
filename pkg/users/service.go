package users

import (
	"context"

	"github.com/Tra-Dew/users/pkg/core"
	"github.com/google/uuid"
)

type service struct {
	repository Repository
}

// NewService ...
func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

// Create ...
func (s *service) Create(ctx context.Context, correlationID string, req *CreateUserRequest) (*CreateUserResponse, error) {

	password, err := NewPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user, err := NewUser(uuid.NewString(), req.Name, req.Email, password)
	if err != nil {
		return nil, err
	}

	if err := s.repository.Insert(ctx, user); err != nil {
		return nil, err
	}

	return &CreateUserResponse{ID: user.ID}, nil
}

// Login ...
func (s *service) Login(ctx context.Context, correlationID string, req *LoginRequest) (*LoginResponse, error) {

	user, err := s.repository.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, core.ErrInvalidCredentials
	}

	if err = user.Password.Equal(req.Password); err != nil {
		return nil, core.ErrInvalidCredentials
	}

	token := "<token>"
	return &LoginResponse{Token: token}, nil
}

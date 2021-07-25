package users

import (
	"context"
	"time"

	"github.com/d-leme/tradew-users/pkg/core"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type service struct {
	settings   *core.Settings
	repository Repository
}

// NewService ...
func NewService(settings *core.Settings, repository Repository) Service {
	return &service{
		settings:   settings,
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Duration(s.settings.JWT.ExpirationInMinutes) * time.Minute),
	})

	tokenString, err := token.SignedString([]byte(s.settings.JWT.Secret))
	if err != nil {
		return nil, err
	}

	return &LoginResponse{Token: tokenString}, nil
}

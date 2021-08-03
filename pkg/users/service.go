package users

import (
	"context"
	"time"

	"github.com/d-leme/tradew-users/pkg/core"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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

	fields := logrus.Fields{
		"email":          req.Email,
		"correlation_id": correlationID,
	}

	password, err := NewPassword(req.Password)
	if err != nil {
		logrus.WithError(err).WithFields(fields).Error("password validation failed")
		return nil, err
	}

	user, err := NewUser(uuid.NewString(), req.Name, req.Email, password)
	if err != nil {
		logrus.WithError(err).WithFields(fields).Error("user validation failed")
		return nil, err
	}

	if err := s.repository.Insert(ctx, user); err != nil {
		logrus.WithError(err).WithFields(fields).Error("error while inserting user")
		return nil, err
	}

	logrus.WithFields(fields).Info("user created successfully")

	return &CreateUserResponse{ID: user.ID}, nil
}

// Login ...
func (s *service) Login(ctx context.Context, correlationID string, req *LoginRequest) (*LoginResponse, error) {

	fields := logrus.Fields{
		"email":          req.Email,
		"correlation_id": correlationID,
	}

	user, err := s.repository.GetByEmail(ctx, req.Email)
	if err != nil {
		logrus.WithError(err).WithFields(fields).Error("error while getting user")
		return nil, err
	}

	if user == nil {
		return nil, core.ErrInvalidCredentials
	}

	if err = user.Password.Equal(req.Password); err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Duration(s.settings.JWT.ExpirationInMinutes) * time.Minute),
	})

	tokenString, err := token.SignedString([]byte(s.settings.JWT.Secret))
	if err != nil {
		logrus.WithError(err).WithFields(fields).Error("error while generating token")
		return nil, err
	}

	logrus.WithFields(fields).Error("logged in successfully")

	return &LoginResponse{Token: tokenString}, nil
}

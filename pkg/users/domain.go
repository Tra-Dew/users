package users

import (
	"context"
	"net/mail"
	"strings"
	"time"

	"github.com/Tra-Dew/users/pkg/core"
	"golang.org/x/crypto/bcrypt"
)

// Password ...
type Password struct {
	Hash string `bson:"hash"`
}

// User ...
type User struct {
	ID        string    `bson:"id"`
	Name      string    `bson:"name"`
	Email     string    `bson:"email"`
	Password  *Password `bson:"password"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

// Service ...
type Service interface {
	Create(ctx context.Context, correlationID string, req *CreateUserRequest) (*CreateUserResponse, error)
	Login(ctx context.Context, correlationID string, req *LoginRequest) (*LoginResponse, error)
}

// Repository ...
type Repository interface {
	Insert(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
}

// NewUser ...
func NewUser(id, name, email string, password *Password) (*User, error) {

	if id == "" {
		return nil, core.ErrValidationFailed
	}

	fixName := strings.TrimSpace(name)
	if fixName == "" {
		return nil, core.ErrValidationFailed
	}

	fixEmail := strings.TrimSpace(email)

	if _, err := mail.ParseAddress(fixEmail); err != nil {
		return nil, core.ErrValidationFailed
	}

	if password == nil {
		return nil, core.ErrValidationFailed
	}

	return &User{
		ID:        id,
		Name:      fixName,
		Email:     fixEmail,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// NewPassword ...
func NewPassword(value string) (*Password, error) {
	if len(value) < 6 {
		return nil, core.ErrValidationFailed
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)

	if err != nil {
		return nil, core.ErrValidationFailed
	}

	return &Password{Hash: string(pass)}, nil
}

// Equal ...
func (currPass *Password) Equal(incomingPass string) error {

	if incomingPass == "" {
		return core.ErrInvalidCredentials
	}

	hash := []byte(currPass.Hash)
	bytePass := []byte(incomingPass)

	err := bcrypt.CompareHashAndPassword(hash, bytePass)
	if err != nil {
		return core.ErrInvalidCredentials
	}

	return nil
}

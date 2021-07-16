package users_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Tra-Dew/users/pkg/users"
	"github.com/Tra-Dew/users/pkg/users/mock"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type serviceTestSuite struct {
	suite.Suite
	assert     *assert.Assertions
	repository *mock.RepositoryMock
	service    users.Service
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}

func (s *serviceTestSuite) SetupSuite() {
	s.assert = assert.New(s.T())
}

func (s *serviceTestSuite) SetupTest() {
	s.repository = mock.NewRepository().(*mock.RepositoryMock)
	s.service = users.NewService(s.repository)
}

func (s *serviceTestSuite) TestCreate() {

	s.repository.On("Insert").Return(nil)

	correlationID := uuid.NewString()
	req := &users.CreateUserRequest{
		Name:     faker.Name(),
		Email:    faker.Email(),
		Password: faker.Password(),
	}

	r, err := s.service.Create(context.TODO(), correlationID, req)

	s.assert.NoError(err)
	s.assert.NotNil(r)
	s.assert.NotEmpty(r.ID)
	s.repository.AssertNumberOfCalls(s.T(), "Insert", 1)
}

func (s *serviceTestSuite) TestCreateInsertError() {

	s.repository.On("Insert").Return(errors.New("unexpected error"))

	correlationID := uuid.NewString()
	req := &users.CreateUserRequest{
		Name:     faker.Name(),
		Email:    faker.Email(),
		Password: faker.Password(),
	}

	r, err := s.service.Create(context.TODO(), correlationID, req)

	s.assert.Error(err)
	s.assert.Nil(r)
	s.repository.AssertNumberOfCalls(s.T(), "Insert", 1)
}

func (s *serviceTestSuite) TestLogin() {
	correlationID := uuid.NewString()
	email := faker.Email()
	value := faker.Password()

	pass, err := users.NewPassword(value)
	s.assert.NoError(err)

	s.repository.On("GetByEmail", email).Return(&users.User{
		ID:        uuid.NewString(),
		Name:      faker.Name(),
		Email:     email,
		Password:  pass,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	r, err := s.service.Login(context.TODO(), correlationID, &users.LoginRequest{
		Email:    email,
		Password: value,
	})

	s.assert.NoError(err)
	s.assert.NotEmpty(r)
	s.repository.AssertNumberOfCalls(s.T(), "GetByEmail", 1)
}

package users_test

import (
	"testing"

	"github.com/Tra-Dew/users/pkg/core"
	"github.com/Tra-Dew/users/pkg/users"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type domainTestSuite struct {
	suite.Suite
	assert *assert.Assertions
}

func TestDomainTestSuite(t *testing.T) {
	suite.Run(t, new(domainTestSuite))
}

func (s *domainTestSuite) SetupSuite() {
	s.assert = assert.New(s.T())
}

func (s *domainTestSuite) TestNewUser() {
	pass, err := users.NewPassword(faker.Password())
	s.assert.NoError(err)
	s.assert.NotNil(pass)

	r, err := users.NewUser(uuid.NewString(), faker.Name(), faker.Email(), pass)

	s.assert.NoError(err)
	s.assert.NotNil(r)
	s.assert.NotEmpty(r.ID)
	s.assert.NotEmpty(r.Name)
	s.assert.NotEmpty(r.Email)
}

func (s *domainTestSuite) TestNewUserInvalidEmail() {

	invalidEmails := []string{
		"invalidemail",
		"invalidemail.com",
		"invalid#email.com",
		"121",
	}

	for _, email := range invalidEmails {
		pass, err := users.NewPassword(faker.Password())
		s.assert.NoError(err)
		s.assert.NotNil(pass)

		r, err := users.NewUser(uuid.NewString(), faker.Name(), email, pass)

		s.assert.Error(err, core.ErrValidationFailed)
		s.assert.Nil(r)
	}
}

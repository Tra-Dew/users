package users_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/d-leme/tradew-users/pkg/core"
	"github.com/d-leme/tradew-users/pkg/users"
	"github.com/d-leme/tradew-users/pkg/users/mock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type controllerTestSuite struct {
	suite.Suite
	assert     *assert.Assertions
	service    *mock.ServiceMock
	controller users.Controller
	engine     *gin.Engine
}

func TestControllerTestSuite(t *testing.T) {
	suite.Run(t, new(controllerTestSuite))
}

func (s *controllerTestSuite) SetupSuite() {
	s.assert = assert.New(s.T())
}

func (s *controllerTestSuite) SetupTest() {
	s.service = mock.NewService().(*mock.ServiceMock)
	s.controller = users.NewController(s.service)

	gin.SetMode(gin.ReleaseMode)
	s.engine = gin.New()
	s.engine.Use(gin.LoggerWithWriter(bytes.NewBufferString("")))
	s.controller.RegisterRoutes(s.engine.Group(""))
}

func (s *controllerTestSuite) TestCreate() {
	id := uuid.NewString()

	req := &users.CreateUserRequest{
		Name:  faker.Name(),
		Email: faker.Email(),
	}

	s.service.On("Create").Return(&users.CreateUserResponse{ID: id}, nil)

	b, err := json.Marshal(req)
	s.assert.NoError(err)

	w := performRequest(s.engine, "POST", "/users", bytes.NewReader(b))
	s.assert.Equal(http.StatusCreated, w.Code)

	var res users.CreateUserResponse
	err = json.NewDecoder(w.Body).Decode(&res)

	h, exists := w.HeaderMap["Location"]

	s.assert.NoError(err)
	s.assert.Equal(id, res.ID)
	s.assert.True(exists)
	s.assert.Len(h, 1)
	s.assert.Equal(h[0], fmt.Sprintf("/users/%s", id))
}

func (s *controllerTestSuite) TestCreateMalformedJson() {
	w := performRequest(s.engine, "POST", "/users", bytes.NewReader([]byte("malformed json")))
	s.assert.Equal(http.StatusUnprocessableEntity, w.Code)

	var res core.RestError
	err := json.NewDecoder(w.Body).Decode(&res)

	s.assert.NoError(err)
	s.assert.Equal(core.ErrMalformedJSON.Key, res.Key)
}

func (s *controllerTestSuite) TestLogin() {
	req := &users.LoginRequest{
		Email:    faker.Email(),
		Password: faker.Password(),
	}

	s.service.On("Login").Return(&users.LoginResponse{Token: "<token>"}, nil)

	b, err := json.Marshal(req)
	s.assert.NoError(err)

	w := performRequest(s.engine, "POST", "/users/login", bytes.NewReader(b))
	s.assert.Equal(http.StatusOK, w.Code)

	var res users.LoginResponse
	err = json.NewDecoder(w.Body).Decode(&res)

	s.assert.NoError(err)
	s.assert.NotEmpty(res.Token)
}

func performRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

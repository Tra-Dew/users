package cmd

import (
	"github.com/Tra-Dew/users/pkg/core"
	"github.com/Tra-Dew/users/pkg/users"
	"github.com/Tra-Dew/users/pkg/users/memory"
)

// Container contains all depencies from our api
type Container struct {
	Settings *core.Settings

	UserRepository users.Repository
	UserService    users.Service
	UserController users.Controller
}

// NewContainer creates new instace of Container
func NewContainer(settings *core.Settings) *Container {

	container := new(Container)

	container.Settings = settings

	container.UserRepository = memory.NewRepository()
	container.UserService = users.NewService(settings, container.UserRepository)
	container.UserController = users.NewController(container.UserService)

	return container
}

// Controllers maps all routes and exposes them
func (c *Container) Controllers() []core.Controller {
	return []core.Controller{
		&c.UserController,
	}
}

// Close terminates every opened resource
func (c *Container) Close() {}

package users

import (
	"fmt"
	"net/http"

	"github.com/Tra-Dew/users/pkg/core"
	"github.com/gin-gonic/gin"
)

// Controller ...
type Controller struct {
	service Service
}

// NewController ...
func NewController(service Service) Controller {
	return Controller{
		service: service,
	}
}

// RegisterRoutes ...
func (c *Controller) RegisterRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.POST("", c.post)
		users.POST("login", c.login)
	}
}

func (c *Controller) post(ctx *gin.Context) {
	req := new(CreateUserRequest)
	correlationID := ctx.GetString("X-Correlation-ID")

	if err := ctx.ShouldBindJSON(req); err != nil {
		core.HandleRestError(ctx, core.ErrMalformedJSON)
		return
	}

	res, err := c.service.Create(ctx, correlationID, req)

	if err != nil {
		core.HandleRestError(ctx, err)
		return
	}

	ctx.Writer.Header().Set("Location", fmt.Sprintf("%s/%s", ctx.Request.URL.Path, res.ID))

	ctx.JSON(http.StatusCreated, res)
}

func (c *Controller) login(ctx *gin.Context) {
	req := new(LoginRequest)
	correlationID := ctx.GetString("X-Correlation-ID")

	if err := ctx.ShouldBindJSON(req); err != nil {
		core.HandleRestError(ctx, core.ErrMalformedJSON)
		return
	}

	res, err := c.service.Login(ctx, correlationID, req)

	if err != nil {
		core.HandleRestError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

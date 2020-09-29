package http

import (
	"github.com/gin-gonic/gin"
	"github.com/micro-gis/oauth-api/src/domain/access_token"
	errors "github.com/micro-gis/oauth-api/src/utils/errors_util"
	"net/http"
)

type AccessTokenHandler interface {
	GetById(c *gin.Context)
	Create(c *gin.Context)
}

type accessTokenHandler struct {
	service access_token.Service
}

func NewHandler(service access_token.Service) AccessTokenHandler {
	return &accessTokenHandler{
		service: service,
	}
}

func (handler *accessTokenHandler) GetById(c *gin.Context) {
	accessTokenId := c.Param("access_token_id")
	accessToken, err := handler.service.GetById(accessTokenId)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, accessToken)
}

func (handler *accessTokenHandler) Create(c *gin.Context) {
	var at access_token.AccessToken
	if err := c.ShouldBindJSON(&at); err != nil {
		restErr := errors.NewBadRequestError("invalid Json Body")
		c.JSON(restErr.Status, restErr)
		return
	}
	if err := handler.service.Create(at); err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusCreated, at)
}

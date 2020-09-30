package app

import (
	"github.com/gin-gonic/gin"
	"github.com/micro-gis/oauth-api/src/clients/cassandra"
	"github.com/micro-gis/oauth-api/src/http"
	"github.com/micro-gis/oauth-api/src/repository/db"
	"github.com/micro-gis/oauth-api/src/repository/rest"
	atService "github.com/micro-gis/oauth-api/src/service/access_token"
)

var (
	router = gin.Default()
)

func StartApplication() {
	session := cassandra.GetSession()
	defer session.Close()

	atHandler := http.NewHandler(
		atService.NewService(rest.NewRestUsersRepository(), db.NewRepository()))

	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	router.POST("/oauth/access_token", atHandler.Create)
	router.Run(":8087")
}

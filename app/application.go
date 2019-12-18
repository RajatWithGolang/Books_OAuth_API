package app

import (
	"github.com/RajatWithGolang/Books_OAuth_API/http"
	"github.com/RajatWithGolang/Books_OAuth_API/repository/db"
	"github.com/RajatWithGolang/Books_OAuth_API/repository/rest"
	"github.com/RajatWithGolang/Books_OAuth_API/services/access_token"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	atHandler := http.NewAccessTokenHandler(
		access_token.NewService(rest.NewRestUsersRepository(), db.NewRepository()))

	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	router.POST("/oauth/access_token", atHandler.Create)

	router.Run(":8080")
}

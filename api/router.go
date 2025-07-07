package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func InitRouter(handler *Handler) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/account", ginify(handler.GetAccount))
		api.POST("/account", ginify(handler.CreateAccount))
		api.POST("/topup", ginify(handler.TopUp))
		api.POST("/blacklist", ginify(handler.BlacklistAccount))
		api.POST("/transfer", ginify(handler.Transfer))
	}

	return router
}

func ginify(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(c.Writer, c.Request)
	}
}
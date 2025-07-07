package api

import (
	"github.com/gin-gonic/gin"
)

func InitRouter(handler *Handler) *gin.Engine {
	router := gin.New()

	router.GET("/", handler.GetRoot)

	// curl -X GET http://127.0.0.1:8080/user/list -H "X-USER: John"
	router.GET("/user/list", handler.GetListUsers)

	// curl -X GET http://127.0.0.1:8080/user/2 -H "X-USER: John"
	router.GET("/user/:id", handler.GetUser)

	// curl -X POST http://127.0.0.1:8080/user -H "Content-Type: application/json" -d "{\"name\": \"Andy\", \"age\": 22}" -H "X-USER: John"
	router.POST("/user", handler.PostRoot)

	// curl -X DELETE http://127.0.0.1:8080/user/4 -H "X-USER: John"
	router.DELETE("user/:id", handler.DeleteUser)

	// curl -X PUT http://127.0.0.1:8080/user/2 -H "Content-Type: application/json" -d "{\"name\": \"Lalala\", \"age\": 22}" -H "X-USER: John"
	router.PUT("/user/:id", handler.UpdateUserData)

	return router
}
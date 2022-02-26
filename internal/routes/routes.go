package routes

import (
	"api-backend-template/internal/ping"
	"api-backend-template/internal/todos"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	v1 := r.Group("v1/")

	v1.GET("/ping", ping.HandlePing)
	v1.GET(todos.RouteGetItems, todos.HandleGetItems)
	v1.POST(todos.RouteCreateItem, todos.HandleCreateItem)
	v1.GET(todos.RouteGetItem, todos.HandleGetItem)
	v1.PUT(todos.RouteUpdateItem, todos.HandleUpdateItem)
	v1.DELETE(todos.RouteDeleteItem, todos.HandleDeleteItem)
}

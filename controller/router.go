package controller
import (
	"app/MQdemo/api"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	router :=gin.Default()
	router.NoRoute(api.NotFound)
	v1 := router.Group("/shop")
	{
		v1.POST("/order",api.Order)
	}
	router.Run(":80")
}

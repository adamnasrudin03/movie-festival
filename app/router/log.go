package router

import (
	"adamnasrudin03/movie-festival/app/controller"
	"adamnasrudin03/movie-festival/app/middlewares"

	"github.com/gin-gonic/gin"
)

func (r routes) logRouter(rg *gin.RouterGroup, logController controller.LogController) {
	logs := rg.Group("/log")
	{
		logs.Use(middlewares.Authentication())
		logs.GET("/", logController.GetAll)
	}
}

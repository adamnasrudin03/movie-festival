package router

import (
	"adamnasrudin03/movie-festival/app/controller"
	"adamnasrudin03/movie-festival/app/middlewares"

	"github.com/gin-gonic/gin"
)

func (r routes) moviesRouter(rg *gin.RouterGroup, moviesController controller.MovieController) {
	movies := rg.Group("/movies")
	{
		movies.Use(middlewares.Authentication())
		movies.POST("/", middlewares.AuthorizationMustBeAdmin(), moviesController.CreateMovie)
		movies.GET("/:id", moviesController.GetOne)
		movies.PUT("/:id", middlewares.AuthorizationMustBeAdmin(), moviesController.UpdateMovie)
		movies.GET("/", moviesController.GetAll)
	}
}

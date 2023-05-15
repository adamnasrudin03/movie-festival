package router

import (
	"adamnasrudin03/movie-festival/app/controller"
	"adamnasrudin03/movie-festival/app/middlewares"

	"github.com/gin-gonic/gin"
)

func (r routes) genresRouter(rg *gin.RouterGroup, genresController controller.GenreController) {
	genres := rg.Group("/genres")
	{
		genres.Use(middlewares.Authentication())
		genres.POST("/", middlewares.AuthorizationMustBeAdmin(), genresController.CreateGenre)
		genres.GET("/:id", middlewares.AuthorizationMustBeAdmin(), genresController.GetOne)
		genres.PUT("/:id", middlewares.AuthorizationMustBeAdmin(), genresController.UpdateGenre)
		genres.DELETE("/:id", middlewares.AuthorizationMustBeAdmin(), genresController.DeleteGenre)
		genres.GET("/", middlewares.AuthorizationMustBeAdmin(), genresController.GetAll)
	}
}

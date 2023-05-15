package router

import (
	"adamnasrudin03/movie-festival/app/controller"
	"adamnasrudin03/movie-festival/pkg/helpers"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type routes struct {
	router *gin.Engine
}

func NewRoutes(contoller controller.Controllers) routes {
	r := routes{
		router: gin.Default(),
	}

	r.router.Use(gin.Logger())
	r.router.Use(gin.Recovery())
	r.router.Use(cors.Default())

	r.router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, helpers.APIResponse("welcome its server", http.StatusOK, nil))
	})

	v1 := r.router.Group("/api/v1")
	r.userRouter(v1, contoller.User)
	r.userRouterAuth(v1, contoller.User)
	r.logRouter(v1, contoller.Log)
	r.moviesRouter(v1, contoller.Movie)

	r.router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, helpers.APIResponse("page not found", http.StatusNotFound, nil))
	})
	return r
}

func (r routes) Run(addr string) error {
	return r.router.Run(addr)
}

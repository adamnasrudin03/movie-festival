package router

import (
	"adamnasrudin03/movie-festival/app/controller"
	"adamnasrudin03/movie-festival/app/middlewares"

	"github.com/gin-gonic/gin"
)

func (r routes) userRouter(rg *gin.RouterGroup, userController controller.UserController) {
	users := rg.Group("/auth")
	{
		users.POST("/sign-up", userController.RegisterUser)
		users.POST("/sign-in", userController.Login)
	}
}

func (r routes) userRouterAuth(rg *gin.RouterGroup, userController controller.UserController) {
	users := rg.Group("/auth-admin")
	{
		users.Use(middlewares.Authentication())
		users.POST("/sign-up", middlewares.AuthorizationMustBeAdmin(), userController.Register)
	}
}

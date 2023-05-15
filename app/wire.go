package app

import (
	"adamnasrudin03/movie-festival/app/controller"
	"adamnasrudin03/movie-festival/app/repository"
	"adamnasrudin03/movie-festival/app/service"

	"gorm.io/gorm"
)

func WiringRepository(db *gorm.DB) *repository.Repositories {
	return &repository.Repositories{
		User: repository.NewUserRepository(db),
	}
}

func WiringService(repo *repository.Repositories) *service.Services {
	return &service.Services{
		User: service.NewUserService(repo.User),
	}
}

func WiringController(srv *service.Services) *controller.Controllers {
	return &controller.Controllers{
		User: controller.NewUserController(srv.User),
	}
}

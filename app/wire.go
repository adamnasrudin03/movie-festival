package app

import (
	"adamnasrudin03/movie-festival/app/controller"
	"adamnasrudin03/movie-festival/app/repository"
	"adamnasrudin03/movie-festival/app/service"

	"gorm.io/gorm"
)

func WiringRepository(db *gorm.DB) *repository.Repositories {
	return &repository.Repositories{
		User:  repository.NewUserRepository(db),
		Log:   repository.NewLogRepository(db),
		Movie: repository.NewMovieRepository(db),
		Genre: repository.NewGenreRepository(db),
	}
}

func WiringService(repo *repository.Repositories) *service.Services {
	return &service.Services{
		User:  service.NewUserService(repo.User),
		Log:   service.NewLogService(repo.Log),
		Movie: service.NewMovieService(repo.Movie),
		Genre: service.NewGenreService(repo.Genre),
	}
}

func WiringController(srv *service.Services) *controller.Controllers {
	return &controller.Controllers{
		User:  controller.NewUserController(srv.User, srv.Log),
		Log:   controller.NewLogController(srv.Log),
		Movie: controller.NewMovieController(srv.Movie, srv.Log),
		Genre: controller.NewGenreController(srv.Genre, srv.Log),
	}
}

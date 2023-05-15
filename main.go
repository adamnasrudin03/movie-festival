package main

import (
	"fmt"

	"adamnasrudin03/movie-festival/app"
	"adamnasrudin03/movie-festival/app/configs"
	"adamnasrudin03/movie-festival/app/router"
	"adamnasrudin03/movie-festival/pkg/database"

	"gorm.io/gorm"
)

var (
	db          *gorm.DB = database.SetupDbConnection()
	repo                 = app.WiringRepository(db)
	services             = app.WiringService(repo)
	controllers          = app.WiringController(services)
)

func main() {
	defer database.CloseDbConnection(db)
	config := configs.GetInstance()

	r := router.NewRoutes(*controllers)

	listen := fmt.Sprintf(":%v", config.Appconfig.Port)
	r.Run(listen)
}

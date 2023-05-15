package repository

import (
	"adamnasrudin03/movie-festival/app/dto"
	"adamnasrudin03/movie-festival/app/entity"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MovieRepository interface {
	Create(ctx *gin.Context, input entity.Movie) (res entity.Movie, err error)
	GetAll(ctx *gin.Context, queryparam dto.ListParam) (result []entity.Movie, total uint64, err error)
	GetByID(ctx *gin.Context, ID uint64) (result entity.Movie, err error)
	UpdateByID(ctx *gin.Context, ID uint64, input entity.Movie) (result entity.Movie, err error)
	DeleteByID(ctx *gin.Context, ID uint64) (err error)
}

type movieRepo struct {
	DB *gorm.DB
}

func NewMovieRepository(db *gorm.DB) MovieRepository {
	return &movieRepo{
		DB: db,
	}
}

func (repo *movieRepo) Create(ctx *gin.Context, input entity.Movie) (res entity.Movie, err error) {
	if err := repo.DB.Create(&input).Error; err != nil {
		log.Printf("[MovieRepository-Create] error Create new Movie: %+v \n", err)
		return input, err
	}

	return input, err
}

func (repo *movieRepo) GetAll(ctx *gin.Context, queryparam dto.ListParam) (result []entity.Movie, total uint64, err error) {
	offset := queryparam.Limit * (queryparam.Page - 1)
	var totaldata int64
	query := repo.DB.WithContext(ctx)
	if queryparam.Search != "" || len(queryparam.Search) > 0 {
		queryparam.Search = fmt.Sprintf("%" + queryparam.Search + "%")
		query = query.Where("title LIKE  = ? OR description LIKE = ? OR artists LIKE = ? OR genres LIKE = ? ",
			queryparam.Search, queryparam.Search, queryparam.Search, queryparam.Search,
		)
	}

	err = query.Model(&entity.Movie{}).Count(&totaldata).Error
	if err != nil {
		log.Printf("[MovieRepository-GetAll] error count total data: %+v \n", err)
		return
	}

	total = uint64(totaldata)
	err = query.Offset(int(offset)).Limit(int(queryparam.Limit)).Find(&result).Error
	if err != nil {
		log.Printf("[MovieRepository-GetAll] error get data: %+v \n", err)
		return
	}

	return
}

func (repo *movieRepo) GetByID(ctx *gin.Context, ID uint64) (result entity.Movie, err error) {
	query := repo.DB.WithContext(ctx)

	if err = query.Where("id = ?", ID).Take(&result).Error; err != nil {
		log.Printf("[MovieRepository-GetByID][%v] error: %+v \n", ID, err)
		return result, err
	}

	return result, err
}

func (repo *movieRepo) UpdateByID(ctx *gin.Context, ID uint64, input entity.Movie) (result entity.Movie, err error) {
	query := repo.DB.WithContext(ctx)

	err = query.Model(&result).Where("id=?", ID).Updates(entity.Movie(input)).Error
	if err != nil {
		log.Printf("[MovieRepository-UpdateByID][%v] error: %+v \n", ID, err)
		return result, err
	}

	return result, err
}

func (repo *movieRepo) DeleteByID(ctx *gin.Context, ID uint64) (err error) {
	query := repo.DB.WithContext(ctx)

	movie := entity.Movie{}
	if err = query.Where("id = ?", ID).Delete(&movie).Error; err != nil {
		log.Printf("[MovieRepository-DeleteByID][%v] error: %+v \n", ID, err)
		return
	}

	return
}

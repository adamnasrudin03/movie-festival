package repository

import (
	"adamnasrudin03/movie-festival/app/dto"
	"adamnasrudin03/movie-festival/app/entity"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GenreRepository interface {
	Create(ctx *gin.Context, input entity.Genre) (res entity.Genre, err error)
	GetAll(ctx *gin.Context, queryparam dto.ListParam) (result []entity.Genre, total uint64, err error)
	GetByID(ctx *gin.Context, ID uint64) (result entity.Genre, err error)
	UpdateByID(ctx *gin.Context, ID uint64, input entity.Genre) (result entity.Genre, err error)
	DeleteByID(ctx *gin.Context, ID uint64) (err error)
}

type genreRepo struct {
	DB *gorm.DB
}

func NewGenreRepository(db *gorm.DB) GenreRepository {
	return &genreRepo{
		DB: db,
	}
}

func (repo *genreRepo) Create(ctx *gin.Context, input entity.Genre) (res entity.Genre, err error) {
	if err := repo.DB.Create(&input).Error; err != nil {
		log.Printf("[GenreRepository-Create] error Create new Genre: %+v \n", err)
		if pgError := err.(*pgconn.PgError); errors.Is(err, pgError) {
			switch pgError.Code {
			case "23505":
				err = errors.New("duplicated key not allowed")
			}
		}

		return input, err
	}

	return input, err
}

func (repo *genreRepo) GetAll(ctx *gin.Context, queryparam dto.ListParam) (result []entity.Genre, total uint64, err error) {
	offset := queryparam.Limit * (queryparam.Page - 1)
	var totaldata int64
	query := repo.DB.WithContext(ctx)
	if queryparam.Search != "" || len(queryparam.Search) > 0 {
		queryparam.Search = "%" + queryparam.Search + "%"
		query = query.Where("lower(name) LIKE ? ", queryparam.Search)
	}

	err = query.Model(&entity.Genre{}).Count(&totaldata).Error
	if err != nil {
		log.Printf("[GenreRepository-GetAll] error count total data: %+v \n", err)
		return
	}

	total = uint64(totaldata)
	err = query.Offset(int(offset)).Limit(int(queryparam.Limit)).Find(&result).Error
	if err != nil {
		log.Printf("[GenreRepository-GetAll] error get data: %+v \n", err)
		return
	}

	return
}

func (repo *genreRepo) GetByID(ctx *gin.Context, ID uint64) (result entity.Genre, err error) {
	query := repo.DB.WithContext(ctx)

	if err = query.Where("id = ?", ID).Take(&result).Error; err != nil {
		log.Printf("[GenreRepository-GetByID][%v] error: %+v \n", ID, err)
		return result, err
	}

	return result, err
}

func (repo *genreRepo) UpdateByID(ctx *gin.Context, ID uint64, input entity.Genre) (result entity.Genre, err error) {
	query := repo.DB.WithContext(ctx)

	err = query.Clauses(clause.Returning{}).Model(&result).Where("id = ?", ID).Updates(entity.Genre(input)).Error
	if err != nil {
		log.Printf("[GenreRepository-UpdateByID][%v] error: %+v \n", ID, err)
		if pgError := err.(*pgconn.PgError); errors.Is(err, pgError) {
			switch pgError.Code {
			case "23505":
				err = errors.New("duplicated key not allowed")
			}
		}

		return result, err
	}

	return result, err
}

func (repo *genreRepo) DeleteByID(ctx *gin.Context, ID uint64) (err error) {
	query := repo.DB.WithContext(ctx)

	Genre := entity.Genre{}
	if err = query.Where("id = ?", ID).Delete(&Genre).Error; err != nil {
		log.Printf("[GenreRepository-DeleteByID][%v] error: %+v \n", ID, err)
		return
	}

	return
}

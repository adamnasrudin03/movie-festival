package service

import (
	"adamnasrudin03/movie-festival/app/dto"
	"adamnasrudin03/movie-festival/app/entity"
	"adamnasrudin03/movie-festival/app/repository"
	"errors"
	"log"
	"math"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GenreService interface {
	Create(ctx *gin.Context, input dto.CreateUpdateGenre) (result entity.Genre, statusCode int, err error)
	GetAll(ctx *gin.Context, queryparam dto.ListParam) (result dto.ResponseList, statusCode int, err error)
	GetByID(ctx *gin.Context, ID uint64) (result entity.Genre, statusCode int, err error)
	UpdateByID(ctx *gin.Context, ID uint64, input dto.CreateUpdateGenre) (result entity.Genre, statusCode int, err error)
	DeleteByID(ctx *gin.Context, ID uint64) (statusCode int, err error)
}

type genreSrv struct {
	GenreRepository repository.GenreRepository
}

func NewGenreService(GenreRepo repository.GenreRepository) GenreService {
	return &genreSrv{
		GenreRepository: GenreRepo,
	}
}

func (srv *genreSrv) Create(ctx *gin.Context, input dto.CreateUpdateGenre) (result entity.Genre, statusCode int, err error) {
	genre := entity.Genre{
		Name: strings.TrimSpace(input.Name),
	}

	result, err = srv.GenreRepository.Create(ctx, genre)
	if err != nil && err.Error() == "duplicated key not allowed" {
		log.Printf("[GenreService-Create] error create data: %+v \n", err)
		return result, http.StatusConflict, errors.New("duplicate record name genre")
	}

	if err != nil {
		log.Printf("[GenreService-Create] error create data: %+v \n", err.Error())
		return result, http.StatusInternalServerError, err
	}

	return result, http.StatusCreated, nil
}

func (srv *genreSrv) GetAll(ctx *gin.Context, queryparam dto.ListParam) (result dto.ResponseList, statusCode int, err error) {
	result.Limit = queryparam.Limit
	result.Page = queryparam.Page

	result.Data, result.Total, err = srv.GenreRepository.GetAll(ctx, queryparam)
	if err != nil {
		log.Printf("[GenreService-GetAll] error get data repo: %+v \n", err)
		return result, http.StatusInternalServerError, err
	}

	result.LastPage = uint64(math.Ceil(float64(result.Total) / float64(queryparam.Limit)))

	return result, http.StatusOK, nil
}

func (srv *genreSrv) GetByID(ctx *gin.Context, ID uint64) (result entity.Genre, statusCode int, err error) {
	result, err = srv.GenreRepository.GetByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) || result.ID == 0 {
		return result, http.StatusNotFound, err
	}

	if err != nil {
		log.Printf("[GenreService-GetByID] error get data repo: %+v \n", err)
		return result, http.StatusInternalServerError, err
	}

	return result, http.StatusOK, nil
}

func (srv *genreSrv) UpdateByID(ctx *gin.Context, ID uint64, input dto.CreateUpdateGenre) (result entity.Genre, statusCode int, err error) {
	genre := entity.Genre{
		Name: strings.TrimSpace(input.Name),
	}

	temp, err := srv.GenreRepository.GetByID(ctx, ID)
	if err != nil && err.Error() == "duplicated key not allowed" {
		log.Printf("[GenreService-Create] error create data: %+v \n", err)
		return result, http.StatusConflict, errors.New("duplicate record name genre")
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || temp.ID == 0 {
		return result, http.StatusNotFound, err
	}

	if err != nil {
		log.Printf("[GenreService-UpdateByID] error get data repo: %+v \n", err)
		return result, http.StatusInternalServerError, err
	}
	result, err = srv.GenreRepository.UpdateByID(ctx, ID, genre)
	if err != nil {
		log.Printf("[GenreService-UpdateByID] error update data repo: %+v \n", err)
		return result, http.StatusInternalServerError, err
	}

	return result, http.StatusOK, nil
}

func (srv *genreSrv) DeleteByID(ctx *gin.Context, ID uint64) (statusCode int, err error) {
	temp, err := srv.GenreRepository.GetByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) || temp.ID == 0 {
		return http.StatusNotFound, err
	}

	if err != nil {
		log.Printf("[GenreService-DeleteByID] error get data repo: %+v \n", err)
		return http.StatusInternalServerError, err
	}

	err = srv.GenreRepository.DeleteByID(ctx, ID)
	if err != nil {
		log.Printf("[GenreService-DeleteByID] error delete data repo: %+v \n", err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

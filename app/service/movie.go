package service

import (
	"adamnasrudin03/movie-festival/app/dto"
	"adamnasrudin03/movie-festival/app/entity"
	"adamnasrudin03/movie-festival/app/repository"
	"errors"
	"log"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MovieService interface {
	Create(ctx *gin.Context, input dto.CreateUpdateMovie) (result entity.Movie, statusCode int, err error)
	GetAll(ctx *gin.Context, queryparam dto.ListParam) (result dto.ResponseList, statusCode int, err error)
	GetByID(ctx *gin.Context, ID uint64) (result entity.Movie, statusCode int, err error)
	UpdateByID(ctx *gin.Context, ID uint64, input dto.CreateUpdateMovie) (result entity.Movie, statusCode int, err error)
	DeleteByID(ctx *gin.Context, ID uint64) (statusCode int, err error)
}

type movieSrv struct {
	MovieRepository repository.MovieRepository
}

func NewMovieService(MovieRepo repository.MovieRepository) MovieService {
	return &movieSrv{
		MovieRepository: MovieRepo,
	}
}

func (srv *movieSrv) Create(ctx *gin.Context, input dto.CreateUpdateMovie) (result entity.Movie, statusCode int, err error) {
	movie := entity.Movie{
		Title:       input.Title,
		Duration:    input.Duration,
		Description: input.Description,
		WatchUrl:    input.WatchUrl,
		Artists:     input.Artists,
		Genres:      input.Genres,
	}
	result, err = srv.MovieRepository.Create(ctx, movie)
	if err != nil {
		log.Printf("[MovieService-Create] error create data: %+v \n", err)
		return result, http.StatusInternalServerError, err
	}

	return result, http.StatusCreated, nil
}

func (srv *movieSrv) GetAll(ctx *gin.Context, queryparam dto.ListParam) (result dto.ResponseList, statusCode int, err error) {
	result.Limit = queryparam.Limit
	result.Page = queryparam.Page

	result.Data, result.Total, err = srv.MovieRepository.GetAll(ctx, queryparam)
	if err != nil {
		log.Printf("[MovieService-GetAll] error get data repo: %+v \n", err)
		return result, http.StatusInternalServerError, err
	}

	result.LastPage = uint64(math.Ceil(float64(result.Total) / float64(queryparam.Limit)))

	return result, http.StatusOK, nil
}

func (srv *movieSrv) GetByID(ctx *gin.Context, ID uint64) (result entity.Movie, statusCode int, err error) {
	result, err = srv.MovieRepository.GetByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) || result.ID == 0 {
		return result, http.StatusNotFound, err
	}

	if err != nil {
		log.Printf("[MovieService-GetByID] error get data repo: %+v \n", err)
		return result, http.StatusInternalServerError, err
	}

	return result, http.StatusOK, nil
}

func (srv *movieSrv) UpdateByID(ctx *gin.Context, ID uint64, input dto.CreateUpdateMovie) (result entity.Movie, statusCode int, err error) {
	movie := entity.Movie{
		Title:       input.Title,
		Duration:    input.Duration,
		Description: input.Description,
		WatchUrl:    input.WatchUrl,
		Artists:     input.Artists,
		Genres:      input.Genres,
	}

	sm, err := srv.MovieRepository.GetByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) || sm.ID == 0 {
		return result, http.StatusNotFound, err
	}

	if err != nil {
		log.Printf("[MovieService-UpdateByID] error get data repo: %+v \n", err)
		return result, http.StatusInternalServerError, err
	}
	result, err = srv.MovieRepository.UpdateByID(ctx, ID, movie)
	if err != nil {
		log.Printf("[MovieService-UpdateByID] error update data repo: %+v \n", err)
		return result, http.StatusInternalServerError, err
	}

	return result, http.StatusOK, nil
}

func (srv *movieSrv) DeleteByID(ctx *gin.Context, ID uint64) (statusCode int, err error) {
	sm, err := srv.MovieRepository.GetByID(ctx, ID)
	if errors.Is(err, gorm.ErrRecordNotFound) || sm.ID == 0 {
		return http.StatusNotFound, err
	}

	if err != nil {
		log.Printf("[MovieService-DeleteByID] error get data repo: %+v \n", err)
		return http.StatusInternalServerError, err
	}

	err = srv.MovieRepository.DeleteByID(ctx, ID)
	if err != nil {
		log.Printf("[MovieService-DeleteByID] error delete data repo: %+v \n", err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

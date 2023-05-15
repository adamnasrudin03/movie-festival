package service

import (
	"log"
	"math"
	"net/http"
	"strings"

	"adamnasrudin03/movie-festival/app/dto"
	"adamnasrudin03/movie-festival/app/entity"
	"adamnasrudin03/movie-festival/app/repository"

	"github.com/gin-gonic/gin"
)

type LogService interface {
	Create(ctx *gin.Context, input entity.Log) (res entity.Log, statusCode int, err error)
	GetAll(ctx *gin.Context, userID uint64, queryparam dto.ListParam) (result dto.ResponseList, statusCode int, err error)
}

type logSrv struct {
	LogRepository repository.LogRepository
}

func NewLogService(LogRepo repository.LogRepository) LogService {
	return &logSrv{
		LogRepository: LogRepo,
	}
}

func (srv *logSrv) Create(ctx *gin.Context, input entity.Log) (res entity.Log, statusCode int, err error) {
	res, err = srv.LogRepository.Create(ctx, input)
	if err != nil {
		log.Printf("[LogService-Register] error create data: %+v \n", err)
		return res, http.StatusInternalServerError, err
	}

	res.User.Role = strings.ToLower(res.User.Role)

	return res, http.StatusCreated, err
}

func (srv *logSrv) GetAll(ctx *gin.Context, userID uint64, queryparam dto.ListParam) (result dto.ResponseList, statusCode int, err error) {
	result.Limit = queryparam.Limit
	result.Page = queryparam.Page

	result.Data, result.Total, err = srv.LogRepository.GetAll(ctx, userID, queryparam)
	if err != nil {
		log.Printf("[LogService-GetAll] error get all data: %+v \n", err)
		return result, http.StatusInternalServerError, err
	}

	result.LastPage = uint64(math.Ceil(float64(result.Total) / float64(queryparam.Limit)))

	return result, http.StatusOK, nil
}

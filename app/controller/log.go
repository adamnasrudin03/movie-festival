package controller

import (
	"net/http"
	"strconv"

	"adamnasrudin03/movie-festival/app/dto"
	"adamnasrudin03/movie-festival/app/service"
	"adamnasrudin03/movie-festival/pkg/helpers"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LogController interface {
	GetAll(ctx *gin.Context)
}

type logController struct {
	Service service.LogService
}

func NewLogController(srv service.LogService) LogController {
	return &logController{
		Service: srv,
	}
}

func (c *logController) GetAll(ctx *gin.Context) {
	var (
		paramPage  uint64 = 1
		paramLimit uint64 = 10
		err        error
	)
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	if ctx.Query("page") == "" {
		paramPage, err = strconv.ParseUint(ctx.Query("page"), 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, helpers.APIResponse("Qwery param page not found or invalid", http.StatusBadRequest, nil))
			return
		}
	}

	if ctx.Query("limit") != "" {
		paramLimit, err = strconv.ParseUint(ctx.Query("limit"), 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, helpers.APIResponse("Query param limit not found or invalid", http.StatusBadRequest, nil))
			return
		}
	}

	param := dto.ListParam{
		Page:  paramPage,
		Limit: paramLimit,
	}
	res, statusHttp, err := c.Service.GetAll(ctx, uint64(userID), param)
	if err != nil {
		ctx.JSON(statusHttp, helpers.APIResponse(err.Error(), statusHttp, nil))
		return
	}

	ctx.JSON(statusHttp, helpers.APIResponse("Get all", statusHttp, res))
}

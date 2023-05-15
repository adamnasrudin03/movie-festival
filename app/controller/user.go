package controller

import (
	"fmt"
	"net/http"
	"strings"

	"adamnasrudin03/movie-festival/app/dto"
	"adamnasrudin03/movie-festival/app/entity"
	"adamnasrudin03/movie-festival/app/service"
	"adamnasrudin03/movie-festival/pkg/helpers"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type UserController interface {
	RegisterUser(ctx *gin.Context)
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type userController struct {
	Service service.UserService
	Log     service.LogService
}

func NewUserController(srv service.UserService, logSrv service.LogService) UserController {
	return &userController{
		Service: srv,
		Log:     logSrv,
	}
}

func (c *userController) Register(ctx *gin.Context) {
	var (
		input dto.RegisterReq
		log   entity.Log
	)
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))
	validate := validator.New()
	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.APIResponse(err.Error(), http.StatusBadRequest, nil))
		return
	}

	err = validate.Struct(input)
	if err != nil {
		errors := helpers.FormatValidationError(err)

		ctx.JSON(http.StatusBadRequest, helpers.APIResponse(errors, http.StatusBadRequest, nil))
		return
	}

	input.Email = strings.TrimSpace(input.Email)
	res, statusHttp, err := c.Service.Register(ctx, input)
	if err != nil {
		ctx.JSON(statusHttp, helpers.APIResponse(err.Error(), statusHttp, nil))
		return
	}
	log.UserID = uint64(userID)
	log.Action = "Create"
	log.Name = fmt.Sprintf("Register new user, user_id = %v", res.ID)
	log.TableName = "Users"
	log.TableNameID = res.ID
	_, _, _ = c.Log.Create(ctx, log)
	ctx.JSON(statusHttp, helpers.APIResponse("User registered", statusHttp, res))
}

func (c *userController) RegisterUser(ctx *gin.Context) {
	var (
		input dto.RegisterUserReq
	)
	validate := validator.New()
	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.APIResponse(err.Error(), http.StatusBadRequest, nil))
		return
	}

	err = validate.Struct(input)
	if err != nil {
		errors := helpers.FormatValidationError(err)

		ctx.JSON(http.StatusBadRequest, helpers.APIResponse(errors, http.StatusBadRequest, nil))
		return
	}

	input.Email = strings.TrimSpace(input.Email)
	req := dto.RegisterReq(input)
	res, statusHttp, err := c.Service.Register(ctx, req)
	if err != nil {
		ctx.JSON(statusHttp, helpers.APIResponse(err.Error(), statusHttp, nil))
		return
	}

	ctx.JSON(statusHttp, helpers.APIResponse("User registered", statusHttp, res))
}

func (c *userController) Login(ctx *gin.Context) {
	var (
		input dto.LoginReq
	)

	validate := validator.New()
	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.APIResponse(err.Error(), http.StatusBadRequest, nil))
		return
	}

	err = validate.Struct(input)
	if err != nil {
		errors := helpers.FormatValidationError(err)

		ctx.JSON(http.StatusBadRequest, helpers.APIResponse(errors, http.StatusBadRequest, nil))
		return
	}

	res, statusHttp, err := c.Service.Login(ctx, input)
	if err != nil {
		ctx.JSON(statusHttp, helpers.APIResponse(err.Error(), statusHttp, nil))
		return
	}

	ctx.JSON(statusHttp, helpers.APIResponse("User login successfully", statusHttp, res))
}

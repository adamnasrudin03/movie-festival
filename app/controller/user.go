package controller

import (
	"fmt"
	"net/http"
	"strings"

	"adamnasrudin03/movie-festival/app/dto"
	"adamnasrudin03/movie-festival/app/service"
	"adamnasrudin03/movie-festival/pkg/helpers"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserController interface {
	RegisterUser(ctx *gin.Context)
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type userController struct {
	Service service.UserService
}

func NewUserController(srv service.UserService) UserController {
	return &userController{
		Service: srv,
	}
}

func (c *userController) Register(ctx *gin.Context) {
	var (
		input dto.RegisterReq
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
	res, statusHttp, err := c.Service.Register(input)
	if err != nil {
		ctx.JSON(statusHttp, helpers.APIResponse(err.Error(), statusHttp, nil))
		return
	}

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
	fmt.Println("input : ", input)
	fmt.Println("req : ", req)
	res, statusHttp, err := c.Service.Register(req)
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

	res, statusHttp, err := c.Service.Login(input)
	if err != nil {
		ctx.JSON(statusHttp, helpers.APIResponse(err.Error(), statusHttp, nil))
		return
	}

	ctx.JSON(statusHttp, helpers.APIResponse("User login successfully", statusHttp, res))
}

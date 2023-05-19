package controller

import (
	"adamnasrudin03/movie-festival/app/dto"
	"adamnasrudin03/movie-festival/app/entity"
	"adamnasrudin03/movie-festival/app/service"
	"adamnasrudin03/movie-festival/pkg/helpers"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type MovieController interface {
	UploadFile(ctx *gin.Context)
	CreateMovie(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	GetOne(ctx *gin.Context)
	UpdateMovie(ctx *gin.Context)
}

type MovieHandler struct {
	Service service.MovieService
	Log     service.LogService
}

func NewMovieController(srv service.MovieService, logSrv service.LogService) MovieController {
	return &MovieHandler{
		Service: srv,
		Log:     logSrv,
	}
}

func (c *MovieHandler) UploadFile(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.APIResponse("No file is received, "+err.Error(), http.StatusBadRequest, nil))
		return
	}

	content := header.Header["Content-Type"]
	if content[0] != "video/mp4" {
		ctx.JSON(http.StatusBadRequest, helpers.APIResponse("Must be file type video", http.StatusBadRequest, nil))
		return
	}

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))
	fileExt := filepath.Ext(header.Filename)
	originalFileName := strings.TrimSuffix(filepath.Base(header.Filename), filepath.Ext(header.Filename))
	now := time.Now()
	filename := fmt.Sprintf("%v-", userID) + strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt

	// config := configs.GetInstance()
	// filePath := "http://localhost:" + config.Appconfig.Port + "/file/uploaded/" + filename
	filePath := "/file/uploaded/" + filename

	out, err := os.Create("public/uploaded/" + filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.APIResponse("Unable to save the file, "+err.Error(), http.StatusInternalServerError, nil))
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.APIResponse("Unable to save the file, "+err.Error(), http.StatusInternalServerError, nil))
		return
	}

	// Logging user upload file
	logging := entity.Log{
		UserID:    uint64(userID),
		Action:    "Create",
		Name:      fmt.Sprintf("Create upload file : %v", filename),
		TableName: "Upload-File",
	}

	go func(ctx *gin.Context, logging entity.Log) {
		_, _, _ = c.Log.Create(ctx, logging)
	}(ctx, logging)

	ctx.JSON(http.StatusCreated, helpers.APIResponse("Your file has been successfully uploaded.", http.StatusCreated, dto.FileRes{Path: filePath}))
}
func (c *MovieHandler) CreateMovie(ctx *gin.Context) {
	var (
		input   dto.CreateUpdateMovie
		logging entity.Log
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

	res, statusHttp, err := c.Service.Create(ctx, input)
	if err != nil {
		ctx.JSON(statusHttp, helpers.APIResponse(err.Error(), statusHttp, nil))
		return
	}

	logging.UserID = uint64(userID)
	logging.Action = "Create"
	logging.Name = fmt.Sprintf("Create new movie, id = %v", res.ID)
	logging.TableName = "Movies"
	logging.TableNameID = res.ID
	go func(ctx *gin.Context, logging entity.Log) {
		_, _, _ = c.Log.Create(ctx, logging)
	}(ctx, logging)

	ctx.JSON(statusHttp, helpers.APIResponse("Movie created", statusHttp, res))
}

func (c *MovieHandler) GetAll(ctx *gin.Context) {
	var (
		paramPage  uint64 = 1
		paramLimit uint64 = 10
		err        error
	)

	if ctx.Query("page") == "" {
		paramPage, err = strconv.ParseUint(ctx.Query("page"), 10, 32)
		if err != nil {
			err = errors.New("query param page invalid")
			ctx.JSON(http.StatusBadRequest, helpers.APIResponse(err.Error(), http.StatusBadRequest, nil))
			return
		}
	}

	if ctx.Query("limit") != "" {
		paramLimit, err = strconv.ParseUint(ctx.Query("limit"), 10, 32)
		if err != nil {
			err = errors.New("query param limit invalid")
			ctx.JSON(http.StatusBadRequest, helpers.APIResponse(err.Error(), http.StatusBadRequest, nil))
			return
		}
	}

	param := dto.ListParam{
		Page:   paramPage,
		Limit:  paramLimit,
		Search: strings.ToLower(strings.TrimSpace(ctx.Query("search"))),
	}

	res, statusHttp, err := c.Service.GetAll(ctx, param)
	if err != nil {
		ctx.JSON(statusHttp, helpers.APIResponse(err.Error(), statusHttp, nil))
		return
	}

	ctx.JSON(statusHttp, helpers.APIResponse("Get all success", statusHttp, res))
}

func (c *MovieHandler) GetOne(ctx *gin.Context) {
	var (
		logging entity.Log
	)
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	ID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		err = errors.New("invalid parameter id")
		ctx.JSON(http.StatusBadRequest, helpers.APIResponse(err.Error(), http.StatusBadRequest, nil))
		return
	}

	res, statusHttp, err := c.Service.GetByID(ctx, ID)
	if err != nil {
		ctx.JSON(statusHttp, helpers.APIResponse(err.Error(), statusHttp, nil))
		return
	}

	logging.UserID = uint64(userID)
	logging.Action = "Get Detail"
	logging.Name = fmt.Sprintf("Get detail movie by id = %v", res.ID)
	logging.TableName = "Movies"
	logging.TableNameID = res.ID
	go func(ctx *gin.Context, logging entity.Log) {
		_, _, _ = c.Log.Create(ctx, logging)
	}(ctx, logging)

	ctx.JSON(statusHttp, helpers.APIResponse("Get Detail movie", statusHttp, res))
}

func (c *MovieHandler) UpdateMovie(ctx *gin.Context) {
	var (
		input   dto.CreateUpdateMovie
		logging entity.Log
	)
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	ID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		err = errors.New("invalid parameter id")
		ctx.JSON(http.StatusBadRequest, helpers.APIResponse(err.Error(), http.StatusBadRequest, nil))
		return
	}

	validate := validator.New()
	err = ctx.ShouldBindJSON(&input)
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

	res, statusHttp, err := c.Service.UpdateByID(ctx, ID, input)
	if err != nil {
		ctx.JSON(statusHttp, helpers.APIResponse(err.Error(), statusHttp, nil))
		return
	}

	logging.UserID = uint64(userID)
	logging.Action = "Update"
	logging.Name = fmt.Sprintf("Update movie by id = %v", res.ID)
	logging.TableName = "Movies"
	logging.TableNameID = res.ID
	go func(ctx *gin.Context, logging entity.Log) {
		_, _, _ = c.Log.Create(ctx, logging)
	}(ctx, logging)

	ctx.JSON(statusHttp, helpers.APIResponse("Movie updated", statusHttp, res))
}

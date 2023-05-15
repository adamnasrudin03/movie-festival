package repository

import (
	"log"

	"adamnasrudin03/movie-festival/app/dto"
	"adamnasrudin03/movie-festival/app/entity"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LogRepository interface {
	Create(ctx *gin.Context, input entity.Log) (res entity.Log, err error)
	GetAll(ctx *gin.Context, userID uint64, queryparam dto.ListParam) (result []entity.Log, total uint64, err error)
}

type logRepo struct {
	DB *gorm.DB
}

func NewLogRepository(db *gorm.DB) LogRepository {
	return &logRepo{
		DB: db,
	}
}

func (repo *logRepo) Create(ctx *gin.Context, input entity.Log) (res entity.Log, err error) {
	if err := repo.DB.Create(&input).Error; err != nil {
		log.Printf("[LogRepository-Create] error Create new Log: %+v \n", err)
		return input, err
	}

	return input, err
}

func (repo *logRepo) GetAll(ctx *gin.Context, userID uint64, queryparam dto.ListParam) (result []entity.Log, total uint64, err error) {
	offset := queryparam.Limit * (queryparam.Page - 1)
	var totaldata int64
	query := repo.DB.WithContext(ctx)
	err = query.Where("user_id = ? ", userID).Model(&entity.Log{}).Count(&totaldata).Error
	if err != nil {
		log.Printf("[LogRepository-GetAll] error count total data: %+v \n", err)
		return
	}
	total = uint64(totaldata)
	err = query.Where("user_id = ? ", userID).Offset(int(offset)).Limit(int(queryparam.Limit)).Preload(clause.Associations).
		Order("updated_at desc").Find(&result).Error
	if err != nil {
		log.Printf("[LogRepository-GetAll] error get data: %+v \n", err)
		return
	}

	return
}

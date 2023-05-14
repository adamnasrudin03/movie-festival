package repository

import (
	"errors"
	"log"

	"adamnasrudin03/movie-festival/app/dto"
	"adamnasrudin03/movie-festival/app/entity"
	"adamnasrudin03/movie-festival/pkg/helpers"

	"gorm.io/gorm"
)

type UserRepository interface {
	Register(input entity.User) (res entity.User, err error)
	Login(input dto.LoginReq) (res entity.User, er error)
	GetByEmail(email string) (res entity.User, err error)
}

type userRepo struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{
		DB: db,
	}
}

func (repo *userRepo) Register(input entity.User) (res entity.User, err error) {
	if err := repo.DB.Create(&input).Error; err != nil {
		log.Printf("[UserRepository-Register] error register new user: %+v \n", err)
		return input, err
	}

	return input, err
}

func (repo *userRepo) Login(input dto.LoginReq) (res entity.User, err error) {
	if err = repo.DB.Where("email = ?", input.Email).Take(&res).Error; err != nil {
		log.Printf("[UserRepository-Login] error login: %+v \n", err)
		return
	}

	if !helpers.PasswordValid(res.Password, input.Password) {
		err = errors.New("invalid password")
		log.Printf("[UserRepository-Login] error cek pass: %+v \n", err)
		return
	}
	return
}

func (repo *userRepo) GetByEmail(email string) (res entity.User, err error) {
	if err = repo.DB.Where("email = ?", email).Take(&res).Error; err != nil {
		log.Printf("[UserRepository-GetByEmail] error : %+v \n", err)
		return
	}
	return
}

package repository

import (
	"adamnasrudin03/movie-festival/app/dto"
	"adamnasrudin03/movie-festival/app/entity"
	"errors"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MovieRepository interface {
	Create(ctx *gin.Context, input entity.Movie, inputGenres []entity.GenreMovies) (res dto.MovieRes, err error)
	GetAll(ctx *gin.Context, queryparam dto.ListParam) (result []entity.Movie, total uint64, err error)
	GetByID(ctx *gin.Context, ID uint64) (result entity.Movie, err error)
	UpdateByID(ctx *gin.Context, ID uint64, input entity.Movie, inputGenres []entity.GenreMovies) (result dto.MovieRes, err error)
	DeleteByID(ctx *gin.Context, ID uint64) (err error)
}

type movieRepo struct {
	DB *gorm.DB
}

func NewMovieRepository(db *gorm.DB) MovieRepository {
	return &movieRepo{
		DB: db,
	}
}

func (repo *movieRepo) Create(ctx *gin.Context, input entity.Movie, inputGenres []entity.GenreMovies) (res dto.MovieRes, err error) {
	query := repo.DB.WithContext(ctx)
	query = query.Begin()

	var genresArray []string
	// check to table genres
	for _, v := range inputGenres {
		genre := entity.Genre{}
		if err = query.Where("id = ?", v.GenreID).Take(&genre).Error; err != nil {
			log.Printf("[MovieRepository-Create][%v] error check record genre: %+v \n", v.GenreID, err)
			query.Rollback()
			return res, err
		}
		genresArray = append(genresArray, genre.Name)
	}

	input.Genres = strings.Join(genresArray, ", ")

	// Create to table movies
	if err := query.Clauses(clause.Returning{}).Create(&input).Error; err != nil {
		log.Printf("[MovieRepository-Create] error Create new Movie: %+v \n", err)
		query.Rollback()
		return res, err
	}

	// append data movie_id in genre_movies
	genreMovies := []entity.GenreMovies{}
	for _, v := range inputGenres {
		genreMovies = append(genreMovies, entity.GenreMovies{GenreID: v.GenreID, MovieID: input.ID})
	}

	if len(genreMovies) > 0 {
		// create to table genre_movies
		if err = repo.DB.Create(&genreMovies).Error; err != nil {
			log.Printf("[MovieRepository-Create] error Create new genre movies: %+v \n", err)
			query.Rollback()

			if pgError := err.(*pgconn.PgError); errors.Is(err, pgError) {
				switch pgError.Code {
				case "23505":
					err = errors.New("duplicated key not allowed")
				}
			}
			return res, err
		}
	}

	query.Commit()
	res = dto.MovieRes{
		ID:           input.ID,
		Title:        input.Title,
		Description:  input.Description,
		Duration:     input.Duration,
		DurationType: input.DurationType,
		WatchUrl:     input.WatchUrl,
		Artists:      input.Artists,
		Genres:       input.Genres,
		Viewers:      input.Viewers,
		CreatedAt:    input.CreatedAt,
		UpdatedAt:    input.UpdatedAt,
	}

	return res, err
}

func (repo *movieRepo) GetAll(ctx *gin.Context, queryparam dto.ListParam) (result []entity.Movie, total uint64, err error) {
	offset := queryparam.Limit * (queryparam.Page - 1)
	var totaldata int64
	query := repo.DB.WithContext(ctx)
	if queryparam.Search != "" || len(queryparam.Search) > 0 {
		queryparam.Search = "%" + queryparam.Search + "%"
		query = query.Where("lower(title) LIKE ? OR lower(description) LIKE ? OR lower(artists) LIKE ? OR lower(genres) LIKE ? ",
			queryparam.Search, queryparam.Search, queryparam.Search, queryparam.Search,
		)
	}

	err = query.Model(&entity.Movie{}).Count(&totaldata).Error
	if err != nil {
		log.Printf("[MovieRepository-GetAll] error count total data: %+v \n", err)
		return
	}

	total = uint64(totaldata)
	err = query.Offset(int(offset)).Limit(int(queryparam.Limit)).Order("viewers desc").Find(&result).Error
	if err != nil {
		log.Printf("[MovieRepository-GetAll] error get data: %+v \n", err)
		return
	}

	return
}

func (repo *movieRepo) GetByID(ctx *gin.Context, ID uint64) (result entity.Movie, err error) {
	query := repo.DB.WithContext(ctx)

	if err = query.Where("id = ?", ID).Take(&result).Error; err != nil {
		log.Printf("[MovieRepository-GetByID][%v] error: %+v \n", ID, err)
		return result, err
	}

	query = query.Begin()
	err = query.Clauses(clause.Returning{}).Model(&result).Where("id = ?", ID).Updates(entity.Movie{Viewers: result.Viewers + 1}).Error
	if err != nil {
		log.Printf("[MovieRepository-GetByID][%v] error update viewers: %+v \n", ID, err)
		query.Rollback()
		return result, err
	}

	if err = query.Commit().Error; err != nil {
		log.Printf("[MovieRepository-GetByID][%v] error: %+v \n", ID, err)
		return result, err
	}

	return result, err
}

func (repo *movieRepo) UpdateByID(ctx *gin.Context, ID uint64, input entity.Movie, inputGenres []entity.GenreMovies) (result dto.MovieRes, err error) {
	query := repo.DB.WithContext(ctx)
	temp := entity.Movie{}
	query = query.Begin()

	if err = query.Where("id = ?", ID).Take(&temp).Error; err != nil {
		log.Printf("[MovieRepository-UpdateByID][%v] error get by id: %+v \n", ID, err)
		return result, err
	}

	err = query.Clauses(clause.Returning{}).Model(&temp).Where("id=?", ID).Updates(entity.Movie(input)).Error
	if err != nil {
		log.Printf("[MovieRepository-UpdateByID][%v] error update record movies: %+v \n", ID, err)
		query.Rollback()
		return result, err
	}

	result = dto.MovieRes{
		ID:           temp.ID,
		Title:        temp.Title,
		Description:  temp.Description,
		Duration:     temp.Duration,
		DurationType: temp.DurationType,
		Artists:      temp.Artists,
		Genres:       temp.Genres,
		Viewers:      temp.Viewers,
		WatchUrl:     temp.WatchUrl,
		CreatedAt:    temp.CreatedAt,
		UpdatedAt:    temp.UpdatedAt,
	}

	// If not update genre movies
	if len(inputGenres) == 0 {
		query.Commit()
		return result, err
	}

	// check to table genres
	var genresArray []string
	var resultGenre []dto.GenreRes
	for _, v := range inputGenres {
		genre := entity.Genre{}
		if err = query.Where("id = ?", v.GenreID).Take(&genre).Error; err != nil {
			log.Printf("[MovieRepository-UpdateByID][%v] error check record genre: %+v \n", v.GenreID, err)
			query.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) || genre.ID == 0 {
				err = errors.New("record genre not found")
			}
			return result, err
		}
		resultGenre = append(resultGenre, dto.GenreRes{ID: v.GenreID, Name: genre.Name})
		genresArray = append(genresArray, genre.Name)
	}

	genreStrings := strings.Join(genresArray, ", ")

	// if input update genre movies not match in record genre movies
	if temp.Genres != genreStrings {
		// Update genres text
		err = query.Clauses(clause.Returning{}).Model(&temp).Where("id=?", ID).Updates(entity.Movie{Genres: genreStrings}).Error
		if err != nil {
			log.Printf("[MovieRepository-UpdateByID][%v] error update genres in table movies: %+v \n", ID, err)
			query.Rollback()
			return result, err
		}

		// Drop record genere_movies
		deleteGM := entity.GenreMovies{}
		if err = query.Where("movie_id = ?", ID).Delete(&deleteGM).Error; err != nil {
			log.Printf("[MovieRepository-UpdateByID][%v]  error delete record genre_movies: %+v \n", ID, err)
			query.Rollback()
			return result, err
		}

		// create new record to table genre_movies
		if err = query.Create(&inputGenres).Error; err != nil {
			log.Printf("[MovieRepository-UpdateByID][%v]  error Create new genre_movies: %+v \n", ID, err)
			query.Rollback()
			if pgError := err.(*pgconn.PgError); errors.Is(err, pgError) {
				switch pgError.Code {
				case "23505":
					err = errors.New("duplicated key not allowed")
				}
			}
			return result, err
		}

		result.GenreDetails = resultGenre
		result.Genres = genreStrings
	}

	query.Commit()
	return result, err
}

func (repo *movieRepo) DeleteByID(ctx *gin.Context, ID uint64) (err error) {
	query := repo.DB.WithContext(ctx)

	movie := entity.Movie{}
	if err = query.Where("id = ?", ID).Delete(&movie).Error; err != nil {
		log.Printf("[MovieRepository-DeleteByID][%v] error: %+v \n", ID, err)
		return
	}

	return
}

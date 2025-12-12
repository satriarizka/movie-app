package usecase

import (
	"errors"
	"math"
	"movie-app/internal/delivery/http/dto/request"
	"movie-app/internal/domain"
	"movie-app/internal/repository"
	"movie-app/pkg/utils"

	"github.com/google/uuid"
)

type MovieUseCase interface {
	Create(req request.CreateMovieRequest) (*domain.Movie, error)
	Update(id uuid.UUID, req request.UpdateMovieRequest) (*domain.Movie, error)
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (*domain.Movie, error)
	GetAll(page int, limit int, search string) ([]domain.Movie, *utils.PaginationMeta, error)
}

type movieUseCase struct {
	movieRepo repository.MovieRepository
}

func NewMovieUseCase(movieRepo repository.MovieRepository) MovieUseCase {
	return &movieUseCase{movieRepo}
}

func (uc *movieUseCase) Create(req request.CreateMovieRequest) (*domain.Movie, error) {
	movie := &domain.Movie{
		Title:       req.Title,
		Description: req.Description,
		Duration:    req.Duration,
		Genre:       req.Genre,
		PosterURL:   req.PosterURL,
	}
	if err := uc.movieRepo.Create(movie); err != nil {
		return nil, err
	}
	return movie, nil
}

func (uc *movieUseCase) Update(id uuid.UUID, req request.UpdateMovieRequest) (*domain.Movie, error) {
	movie, err := uc.movieRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("movie not found")
	}

	if req.Title != "" {
		movie.Title = req.Title
	}
	if req.Description != "" {
		movie.Description = req.Description
	}
	if req.Duration > 0 {
		movie.Duration = req.Duration
	}
	if req.Genre != "" {
		movie.Genre = req.Genre
	}
	if req.PosterURL != "" {
		movie.PosterURL = req.PosterURL
	}

	if err := uc.movieRepo.Update(movie); err != nil {
		return nil, err
	}
	return movie, nil
}

func (uc *movieUseCase) Delete(id uuid.UUID) error {
	_, err := uc.movieRepo.FindByID(id)
	if err != nil {
		return errors.New("movie not found")
	}
	return uc.movieRepo.Delete(id)
}

func (uc *movieUseCase) GetByID(id uuid.UUID) (*domain.Movie, error) {
	return uc.movieRepo.FindByID(id)
}

func (uc *movieUseCase) GetAll(page int, limit int, search string) ([]domain.Movie, *utils.PaginationMeta, error) {
	movies, total, err := uc.movieRepo.FindAll(page, limit, search)
	if err != nil {
		return nil, nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	meta := &utils.PaginationMeta{
		CurrentPage: page,
		TotalPage:   totalPages,
		TotalItems:  total,
		Limit:       limit,
	}
	return movies, meta, nil
}

package shop

import (
    "context"

    "github.com/RomanAgaltsev/avito-shop/internal/config"
    "github.com/RomanAgaltsev/avito-shop/internal/model"
)

// Service is the user service interface.
type Service interface {
    UserRegister(ctx context.Context, user *model.User) error
    UserLogin(ctx context.Context, user *model.User) error
}

// Repository is the user service repository interface.
type Repository interface {
    CreateUser(ctx context.Context, user *model.User) error
    GetUser(ctx context.Context, login string) (*model.User, error)
}

// NewService creates new user service.
func NewService(repository Repository, cfg *config.Config) (Service, error) {
    return &service{
        repository: repository,
        cfg:        cfg,
    }, nil
}

// service is the user service structure.
type service struct {
    repository Repository
    cfg        *config.Config
}

// UserRegister creates new user.
func (s *service) UserRegister(ctx context.Context, user *model.User) error {
    return nil
}

// UserLogin logins existed user.
func (s *service) UserLogin(ctx context.Context, user *model.User) error {
    return nil
}

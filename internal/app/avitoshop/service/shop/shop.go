package shop

import (
    "context"
    "errors"
    "fmt"

    "github.com/RomanAgaltsev/avito-shop/internal/app/avitoshop/service/repository"
    "github.com/RomanAgaltsev/avito-shop/internal/config"
    "github.com/RomanAgaltsev/avito-shop/internal/model"
    "github.com/RomanAgaltsev/avito-shop/internal/pkg/auth"
)

var (
    _ Service    = (*service)(nil)
    _ Repository = (*repository.Repository)(nil)

    ErrUserNameIsAlreadyTaken = fmt.Errorf("user name has already been taken")
    ErrWrongUserNamePassword  = fmt.Errorf("wrong username/password")
)

// Service is the user service interface.
type Service interface {
    UserAuth(ctx context.Context, user model.User) error
    UserBalance(ctx context.Context, user model.User) error
    UserInfo(ctx context.Context, user model.User) (model.Info, error)
    SendCoins(ctx context.Context, fromUser model.User, toUser model.User, amount int) error
    BuyItem(ctx context.Context, user model.User, item model.InventoryItem) error
}

// Repository is the user service repository interface.
type Repository interface {
    CreateUser(ctx context.Context, user model.User) (model.User, error)
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

// UserAuth creates new user.
func (s *service) UserAuth(ctx context.Context, user model.User) error {
    // Replace password with hash
    hash, err := auth.HashPassword(user.Password)
    if err != nil {
        return err
    }
    user.Password = hash

    // Create user in the repository
    userInRepo, err := s.repository.CreateUser(ctx, user)

    // There is a conflict - user name is already exists in the database
    if errors.Is(err, repository.ErrConflict) {
        return ErrUserNameIsAlreadyTaken
    }

    // If user doesn`t exist or password is wrong
    if !auth.CheckPasswordHash(user.Password, userInRepo.Password) {
        return ErrWrongUserNamePassword
    }

    return nil
}

// UserBalance creates new user balance.
func (s *service) UserBalance(ctx context.Context, user model.User) error {
    return nil
}

// SendCoins sends given amount of coins from one user to another.
func (s *service) SendCoins(ctx context.Context, fromUser model.User, toUser model.User, amount int) error {
    return nil
}

// BuyItem buys a given inventory item.
func (s *service) BuyItem(ctx context.Context, user model.User, item model.InventoryItem) error {
    return nil
}

// UserInfo returns user info about coins, inventory and transaction history.
func (s *service) UserInfo(ctx context.Context, user model.User) (model.Info, error) {
    return model.Info{}, nil
}

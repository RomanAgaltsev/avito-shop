package shop

import (
	"context"

	"github.com/RomanAgaltsev/avito-shop/internal/config"
	"github.com/RomanAgaltsev/avito-shop/internal/model"
)

// Service is the user service interface.
type Service interface {
	UserRegister(ctx context.Context, user model.User) error
	UserLogin(ctx context.Context, user model.User) error
	UserBalance(ctx context.Context, user model.User) error
	SendCoins(ctx context.Context, fromUser model.User, toUser model.User, amount int) error
	BuyItem(ctx context.Context, user model.User, item model.InventoryItem) error
	UserInfo(ctx context.Context, user model.User) (model.Info, error)
}

// Repository is the user service repository interface.
type Repository interface {
	CreateUser(ctx context.Context, user model.User) error
	GetUser(ctx context.Context, login string) (model.User, error)
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
func (s *service) UserRegister(ctx context.Context, user model.User) error {
	return nil
}

// UserLogin logins existed user.
func (s *service) UserLogin(ctx context.Context, user model.User) error {
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

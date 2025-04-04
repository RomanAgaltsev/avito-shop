// Package shop implements business logic off the application.
package shop

import (
	"context"
	"errors"
	"fmt"

	"github.com/cenkalti/backoff/v4"

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
	ErrNotEnoughBalance       = fmt.Errorf("not enough coins to send")
	ErrNoSuchItem             = fmt.Errorf("no such item")
	ErrNoSuchUser             = fmt.Errorf("no such user to send coins")
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
	CreateUser(ctx context.Context, bo *backoff.ExponentialBackOff, user model.User) (model.User, error)
	CreateBalance(ctx context.Context, bo *backoff.ExponentialBackOff, user model.User) error
	SendCoins(ctx context.Context, bo *backoff.ExponentialBackOff, fromUser model.User, toUser model.User, amount int) error
	BuyItem(ctx context.Context, bo *backoff.ExponentialBackOff, user model.User, item model.InventoryItem) error
	GetBalance(ctx context.Context, bo *backoff.ExponentialBackOff, user model.User) (int, error)
	GetInventory(ctx context.Context, bo *backoff.ExponentialBackOff, user model.User) ([]model.InventoryItem, error)
	GetHistory(ctx context.Context, bo *backoff.ExponentialBackOff, user model.User) (model.CoinsHistory, error)
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
	password := user.Password
	user.Password = hash

	// Create user in the repository
	userInRepo, err := s.repository.CreateUser(ctx, repository.DefaultBackOff, user)

	// There is a conflict - user name is already exists in the database
	if errors.Is(err, repository.ErrConflict) && !auth.CheckPasswordHash(password, userInRepo.Password) {
		return ErrWrongUserNamePassword
	}

	if err != nil && !errors.Is(err, repository.ErrConflict) {
		return err
	}

	if err == nil {
		return s.repository.CreateBalance(ctx, repository.DefaultBackOff, user)
	}

	return nil
}

// UserBalance creates new user balance.
func (s *service) UserBalance(ctx context.Context, user model.User) error {
	return s.repository.CreateBalance(ctx, repository.DefaultBackOff, user)
}

// SendCoins sends given amount of coins from one user to another.
func (s *service) SendCoins(ctx context.Context, fromUser model.User, toUser model.User, amount int) error {
	err := s.repository.SendCoins(ctx, repository.DefaultBackOff, fromUser, toUser, amount)
	if errors.Is(err, repository.ErrNoData) {
		return ErrNoSuchUser
	}
	if errors.Is(err, repository.ErrNegativeBalance) {
		return ErrNotEnoughBalance
	}

	if err != nil {
		return err
	}

	return nil
}

// BuyItem buys a given inventory item.
func (s *service) BuyItem(ctx context.Context, user model.User, item model.InventoryItem) error {
	err := s.repository.BuyItem(ctx, repository.DefaultBackOff, user, item)
	if errors.Is(err, repository.ErrNoData) {
		return ErrNoSuchItem
	}
	if errors.Is(err, repository.ErrNegativeBalance) {
		return ErrNotEnoughBalance
	}

	if err != nil {
		return err
	}

	return nil
}

// UserInfo returns user info about coins, inventory and transaction history.
func (s *service) UserInfo(ctx context.Context, user model.User) (model.Info, error) {
	coins, err := s.repository.GetBalance(ctx, repository.DefaultBackOff, user)
	if err != nil {
		return model.Info{}, err
	}

	inventory, err := s.repository.GetInventory(ctx, repository.DefaultBackOff, user)
	if err != nil {
		return model.Info{}, err
	}

	history, err := s.repository.GetHistory(ctx, repository.DefaultBackOff, user)
	if err != nil {
		return model.Info{}, err
	}

	return model.Info{
		Coins:        coins,
		Inventory:    inventory,
		CoinsHistory: history,
	}, nil
}

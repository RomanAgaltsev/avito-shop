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
    ErrNotEnoughBalance       = fmt.Errorf("not enough coins to send")
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
    CreateBalance(ctx context.Context, user model.User) error
    SendCoins(ctx context.Context, fromUser model.User, toUser model.User, amount int) error
    BuyItem(ctx context.Context, user model.User, item model.InventoryItem) error
    GetBalance(ctx context.Context, user model.User) (int, error)
    GetInventory(ctx context.Context, user model.User) ([]model.InventoryItem, error)
    GetHistory(ctx context.Context, user model.User) (model.CoinsHistory, error)
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
    return s.repository.CreateBalance(ctx, user)
}

// SendCoins sends given amount of coins from one user to another.
func (s *service) SendCoins(ctx context.Context, fromUser model.User, toUser model.User, amount int) error {
    err := s.repository.SendCoins(ctx, fromUser, toUser, amount)
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
    err := s.repository.BuyItem(ctx, user, item)
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
    coins, err := s.repository.GetBalance(ctx, user)
    if err != nil {
        return model.Info{}, err
    }

    inventory, err := s.repository.GetInventory(ctx, user)
    if err != nil {
        return model.Info{}, err
    }

    history, err := s.repository.GetHistory(ctx, user)
    if err != nil {
        return model.Info{}, err
    }

    // There is a compromise between the number of database accesses and the memory allocations for slices capacity.
    received := make([]model.CoinsReceiving, 0, len(history))
    sent := make([]model.CoinsSending, 0, len(history))

    for _, rec := range history {
        if rec.FromUser != "" {
            received = append(received, model.CoinsReceiving{
                FromUser: rec.FromUser,
                Amount:   int(rec.Amount),
            })
            continue
        }
        if rec.ToUser != "" {
            sent = append(sent, model.CoinsSending{
                ToUser: rec.ToUser,
                Amount: int(rec.Amount),
            })
        }
    }

    coinHistory := model.CoinsHistory{
        Received: received,
        Sent:     sent,
    }

    return model.Info{
        Coins:        coins,
        Inventory:    inventory,
        CoinsHistory: coinHistory,
    }, nil
}

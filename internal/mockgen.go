package internal

//go:generate mockgen --build_flags=--mod=mod -destination=./mock/repository.go -package=mock github.com/RomanAgaltsev/avito-shop/internal/app/avitoshop/service/shop Repository

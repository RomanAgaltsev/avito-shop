[![codecov](https://codecov.io/gh/RomanAgaltsev/avito-shop/graph/badge.svg?token=H51MAKTK3M)](https://codecov.io/gh/RomanAgaltsev/avito-shop)

# **Тестовое задание для стажёра Backend-направления (зимняя волна 2025)**

## Магазин мерча - "avito-shop"

Описание
задачи [тут](https://github.com/avito-tech/tech-internship/blob/main/Tech%20Internships/Backend/Backend-trainee-assignment-winter-2025/Backend-trainee-assignment-winter-2025.md)

* Автор: Роман Агальцев
* Почта: roman-agalcev@yandex.ru

## Описание проекта

Проект выполнен на Go с использованием PostgreSQL в качестве СУБД.

Реализованы следующие хендлеры:

* POST /api/auth - регистрация и аутентификация пользователей
* GET /api/buy/{item} - приобретение пользователем мерча
* POST /api/sendCoin - отправка одним пользователем монет другому пользователю
* GET /api/info - получение информации о пользователе - его баланс монет, приобретенные вещи, а также история транзакций
  с монетами

Для авторизации пользователей используются JWT-токены. При первой авторизации в БД создается новый пользователь и для
него устанавливается стартовый баланс в 1000 монет.

Приложение разделено на слои и состоит из одного сервиса - shop.

В качестве БД выбрана PostgreSQL. Для взаимодействия с БД используется связка sqlc+pgxpool.
Для pgxpool также подготовлен интерфейс и мок для тестов репозитория.

В проекте использованы следующие пакеты и утилиты:

* go-chi/chi - HTTP роутер
* go-chi/jwtauth - работа с JWT-токенами и аутентификация
* go-chi/render - обработка содержания HTTP запросов и ответов
* cenkalti/backoff - retry операции
* sqlc-dev/sqlc - генерация boilerplate-кода для взаимодействия с СУБД
* jackc/pgx - работа с pgx pool
* onsi/ginkgo и onsi/gomega - подготовка тестов
* pressly/goose - миграции
* samber/slog-chi - отдельный slog логер для запросов в роутере chi
* go.uber.org/mock - моки
* go.uber.org/zap - бекенд для логера slog
* golang.org/x/crypto - хэширование паролей пользователей

Кроме этого:

* Подготовлен Makefile
* Подготовлен файлы с OpenAPI спецификацией
* Подготовлен golangci.yml c конфигурацией линтера
* Подготовлен Dockerfile для приложения магазина
* Подготовлены файлы Docker Compose и скрипт инициализации БД PostgreSQL. Для сборки необходимо создать файл .env с
  переменными окружения.

Подготовлены тесты приложения:

* unit-тесты пакетов приложения - по данным CodeCov покрытие составляет около 62%
* E2E-тесты для сценариев покупки мерча и отправки монет пользователями

## Переменные окружения

Для запуска приложения с использованием Docker Compose, необходимо подготовить .env файл со следующими переменными
окружения:

* RUN_ADDRESS - адрес и порт запуска сервиса, например `:8080`
* DATABASE_URI - строка подключения к базе данных Postgres, например
  `postgres://postgres:12345@localhost:5432/avitoshop?sslmode=disable`
* SECRET_KEY - секретный ключ, используемый при аутентификации пользователей, например `secret`

Кроме этого, для инициализации базы данных приложения на Postgres, в файле переменных окружения необходимо дополнительно
определить переменные:

* POSTGRES_USER - пользователь postgres
* POSTGRES_PASSWORD - пароль postgres
* POSTGRES_DB - база данных postgres
* POSTGRES_APP_USER - пользователь базы данных приложения
* POSTGRES_APP_PASS - пароль пользователя базы данных приложения
* POSTGRES_APP_DB - база данных приложения

Также для Pgadmin требуются следующие переменные окружения:

* PGADMIN_DEFAULT_EMAIL - адрес электронной почты
* PGADMIN_DEFAULT_PASSWORD - пароль по умолчанию

В корневой папке проекта присутствует файл `.env` с использованием которого выполнялос тестирование проекта.

Для удобства работы с проектом подготовлен Makefile:

* Для запуска unit-тестов можно воспользоваться командой `make test-unit`
* Для сборки исполняемого файла проекта можно воспользоваться командной `make build`
* Для сборки Docker Compose можно воспользоваться командой `make dc-build`
* Для запуска контейнеров можно воспользоваться командной `make dc-up`

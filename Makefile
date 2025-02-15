pguser = postgres
pguserpass = postgres
pguri = localhost:5432
dbname = avitoshop
secretkey = secret

# ==============================================================================

#  Local commands

.PHONY: gen-mock
gen-mock:	# Generate mocks
	go generate ./...

.PHONY: test
test:	# Execute the unit tests
	go test -count=1 -v -timeout 30s -coverprofile cover.out ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: ginkgo
ginkgo:	# Execute the ginkgo unit tests
	ginkgo ./...

.PHONY: tidy
tidy:	# Cleanup go.mod
	go mod tidy

.PHONY: fmt
fmt:	# Format *.go and *.proto files using gofmt and clang-format
	gofmt -w .
	goimports -local "github.com/RomanAgaltsev/avito-shop" -w .

.PHONY: cover
cover:	# Show the cover report
	go tool cover -html cover.out

.PHONY: update
update:	# Update dependencies as recorded in the go.mod and go.sum files
	go list -m -u all
	go get -u ./...
	go mod tidy

.PHONY: build
build:	# Build application binary
	go build -o cmd/avitoshop/avitoshop cmd/avitoshop/main.go

.PHONY: run
run:	# Run the application
	./cmd/avitoshop/avitoshop -d="postgres://$(pguser):$(pguserpass)@$(pguri)/$(dbname)" -k=$(secretkey)

# ==============================================================================

#  Docker commands

.PHONY: dc-build
dc-build:	# Build docker compose
	docker-compose -f deployments/compose.yaml -f deployments/avitoshop/compose.yaml -f deployments/postgres/compose.yaml --env-file .env build

.PHONY: dc-up
dc-up:	# Build docker compose
	docker-compose -f deployments/compose.yaml -f deployments/avitoshop/compose.yaml -f deployments/postgres/compose.yaml --env-file .env up -d

.PHONY: dc-down
dc-down:	# Build docker compose
	docker-compose -f deployments/compose.yaml -f deployments/avitoshop/compose.yaml -f deployments/postgres/compose.yaml --env-file .env down
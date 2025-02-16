FROM golang:1.23.6
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -o /avitoshop ./cmd/avitoshop/main.go


FROM alpine
WORKDIR /
COPY --from=0 /avitoshop /avitoshop
EXPOSE 8080
ENTRYPOINT ["/avitoshop"]
# Стадия сборки
FROM golang:1.24.2 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# ❗ СТАТИЧЕСКАЯ СБОРКА — вот фокус:
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app ./cmd/main.go

# Финальный образ — суперчистый
FROM scratch

WORKDIR /app
COPY --from=builder /app/app .

# ❗ Тут даже нет bash, только бинарь
CMD ["/app/app"]

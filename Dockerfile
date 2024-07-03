# Сборка приложения
FROM golang:1.20 AS builder

WORKDIR /app

# Копируем и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем все остальные файлы из корня проекта в /app
COPY . .

# Выполняем сборку приложения
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /project-management ./cmd/main.go

# Второй этап: запуск приложения в минимальном образе Alpine
FROM alpine:latest

WORKDIR /app

# Копируем бинарный файл из предыдущего этапа
COPY --from=builder /project-management /project-management

EXPOSE 8080

# Указываем команду для запуска приложения
CMD ["/project-management"]

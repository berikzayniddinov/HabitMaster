# Используем официальный образ Golang
FROM golang:1.23

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей (go.mod и go.sum)
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod tidy

# Копируем всё содержимое проекта в контейнер
COPY . .

# Меняем рабочую директорию на папку main, где находится main.go
WORKDIR /app/main

# Собираем бинарный файл
RUN go build -o main .

# Указываем команду для запуска приложения
CMD ["./main"]

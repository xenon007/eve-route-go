# Сборка фронтенда
FROM node:20-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package.json frontend/yarn.lock ./
RUN yarn install --frozen-lockfile
COPY frontend .
RUN yarn build

# Сборка Go сервера
FROM golang:1.21-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/frontend/dist ./frontend/dist
RUN go build -o eve-route

# Минимальный образ для запуска
FROM alpine:3.18
WORKDIR /root/
COPY --from=build /app/eve-route .
EXPOSE 8080
CMD ["./eve-route"]

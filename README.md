# EVE Route Go

Сервер на Go для планирования маршрутов в EVE Online. Включает пример Capital Jump Planner (расчёт прыжков с радиусом 5 св.лет) и отдаёт встроенный фронтенд.

## Сборка

Требуются Go 1.21 и Node.js 20.

### Фронтенд

```bash
cd frontend
yarn install
NODE_OPTIONS=--openssl-legacy-provider yarn build
cd ..
```

### Бэкенд

```bash
go build -o eve-route
```

## Запуск

```bash
./eve-route
```

Сервер слушает порт `8080` и предоставляет API `/api/capital`. Веб-интерфейс Capital Jump Planner доступен по адресу `http://localhost:8080/#Capital`.

## Переменные окружения

| Переменная            | Описание                                       | Значение по умолчанию       |
| -------------------- | ---------------------------------------------- | --------------------------- |
| `NODE_OPTIONS`       | используется при сборке фронтенда              | `--openssl-legacy-provider` |
| `PORT`               | порт HTTP-сервера                              | `8080`                      |
| `DATABASE_URL`       | строка подключения к БД (postgres/mongodb/sqlite) | `-`                         |
| `OAUTH_CLIENT_ID`    | OAuth2 Client ID                               | `-`                         |
| `OAUTH_CLIENT_SECRET`| OAuth2 Client Secret                           | `-`                         |
| `REDIRECT_URL`       | OAuth2 Redirect URL                            | `-`                         |
| `CSRF_KEY`           | ключ для CSRF защиты                           | `-`                         |
| `API_SECRET`         | секрет для управления Ansiblex API             | `-`                         |
| `SESSION_KEY`        | ключ для cookie-сессий                         | `-`                         |

Если `DATABASE_URL` не задан, используется встроенное в память хранилище.


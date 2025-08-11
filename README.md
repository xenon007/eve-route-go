# EVE Route Go

Сервер на Go для планирования маршрутов в EVE Online. Включает пример Capital Jump Planner и отдаёт встроенный фронтенд.

## Сборка

Требуются Go 1.21 и Node.js 20.

```bash
# сборка фронтенда
cd frontend
yarn install
yarn build
cd ..

# сборка бинарника с вшитым фронтендом
go build -o eve-route
```

## Запуск

```bash
./eve-route
```

Сервер слушает порт `8080` и предоставляет API `/api/capital`. Веб-интерфейс Capital Jump Planner доступен по адресу `http://localhost:8080/#Capital`.

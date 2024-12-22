# Calc Service

Calc Service - это веб-сервис для вычисления арифметических выражений. Пользователь отправляет выражение через HTTP POST-запрос, а сервис возвращает результат вычисления.

## Endpoints

### `POST /api/v1/calculate`

Принимает JSON с выражением и возвращает результат или ошибку.

#### Запрос

```json
{
    "expression": "2+2*2"
}
```

# Инструкция по запуску

1. **Убедиться, что установлен Go** (версии 1.20 или выше).

2. **Склонировать репозиторий**:
```bash
   git clone https://github.com/ladnaaaaaa/calc_service.git
   cd calc_service
```
3. **Запустить сервис**:
```bash
go run ./cmd/calc_service/...
```
4. **Проверить работу** (пример запроса с помощью curl в cmd)
```bash
curl -H "Content-Type: application/json" -X POST http://localhost:8080/api/v1/calculate -d "{\"expression\": \"2+2*2\"}"
```

Ожидаемый ответ:
```json
{
  "result": "6"
}
```

5. **Остановить сервис** нажатием Ctrl + C в терминале.

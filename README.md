# Calc Service

Calc Service - это веб-сервис для распределённых вычислений.

# Инструкция по запуску

1. **Убедиться, что установлен Go** (версии 1.20 или выше).

2. **Склонировать репозиторий**:
```bash
   git clone https://github.com/ladnaaaaaa/calc_service.git
   cd calc_service
```
3. **Запустить оркестратор**:
```bash
go run ./cmd/orchestrator/...
```
4. **Запустить агент (в отдельном окне терминала)**:
```bash
go run ./cmd/agent/...
```
5. По адрессу: http://localhost:8080 в браузере будет доступен веб-интерфейс для работы с сервисом

6. Unit-тесты можно запустить из интерфейса IDE или из командной строки

# Эндпоинты

API сервиса имеет следующие эндпоинты: 
- POST /api/v1/calculate
- GET /api/v1/expressions
- GET /api/v1/expressions/:id
- GET /internal/task
- POST /internal/task

# Переменные окружения

Для того, чтобы задать переменные окружения можно запустить файл "environment.bat" (windows) или задать их в ручную:

```
set TIME_ADDITION_MS=1000
set TIME_MULTIPLICATIONS_MS=2000
set TIME_SUBTRACTION_MS=1000
set TIME_DIVISIONS_MS=3000
```

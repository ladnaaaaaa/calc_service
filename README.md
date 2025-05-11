# Calc Service

Calc Service - это веб-сервис для распределённых вычислений, который позволяет выполнять арифметические выражения в распределённом режиме. Сервис использует архитектуру, состоящую из оркестратора и агентов, для эффективного выполнения задач.

## Ссылка на сайт

Сервис доступен по адресу: [http://localhost:8080](http://localhost:8080)

## Инструкция по запуску

1. **Убедитесь, что установлен Go** (версии 1.20 или выше).

2. **Установите GCC**:
   - Для Windows: Скачайте и установите MinGW-w64 с официального сайта по [инструкции](https://programforyou.ru/poleznoe/kak-ustanovit-gcc-dlya-windows?ysclid=majlp37z7w118007909).

3. **Настройте переменную окружения CGO**:
   - Установите переменную окружения `CGO_ENABLED=1` для включения поддержки CGO:
     ```bash
     set CGO_ENABLED=1
     ```

4. **Установите Protocol Buffers Compiler (protoc)**:
   - Скачайте последнюю версию protoc для Windows: [protoc-25.1-win64.zip](https://github.com/protocolbuffers/protobuf/releases/download/v25.1/protoc-25.1-win64.zip)
   - Распакуйте содержимое архива в папку `tools\protoc` проекта
   - В результате должна появиться структура:
     ```
     tools/
     └── protoc/
         ├── bin/
         │   └── protoc.exe
         └── include/
     ```

5. **Установите Go плагины для protoc**:
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

6. **Сгенерируйте gRPC код**:
   ```bash
   .\scripts\generate_proto.bat
   ```
   После выполнения в папке `api` должны появиться файлы:
   - `calculator.pb.go`
   - `calculator_grpc.pb.go`

7. **Скачайте зависимости проекта**:
   ```bash
   go mod tidy
   ```

8. **Запустите оркестратор**:
   ```bash
   go run ./cmd/orchestrator/...
   ```
   Оркестратор запустит:
   - HTTP-сервер на порту 8080 (для веб-интерфейса)
   - gRPC-сервер на порту 50051 (для общения с агентами)

9. **Запустите агент (в отдельном окне терминала)**:
   ```bash
   go run ./cmd/agent/...
   ```

10. **Откройте веб-интерфейс**:
    - Перейдите по адресу [http://localhost:8080](http://localhost:8080) в браузере.
    - Пройдите авторизацию.
    - Введите арифметическое выражение и нажмите "Вычислить".

## Эндпоинты API

Сервис предоставляет следующие эндпоинты:
- `POST /api/v1/register` - регистрация нового пользователя.
- `POST /api/v1/login` - авторизация пользователя.
- `POST /api/v1/calculate` - отправка выражения на вычисление.
- `GET /api/v1/expressions` - получение списка выражений.
- `GET /api/v1/expressions/:id` - получение выражения по ID.

## gRPC API

Сервис предоставляет следующие gRPC методы:
- `GetTask` - получение задачи для вычисления.
- `SubmitResult` - отправка результата вычисления.

## Переменные окружения

Для настройки времени выполнения операций можно задать следующие переменные окружения:
- `TIME_ADDITION_MS` - время выполнения сложения (мс).
- `TIME_MULTIPLICATIONS_MS` - время выполнения умножения (мс).
- `TIME_SUBTRACTION_MS` - время выполнения вычитания (мс).
- `TIME_DIVISIONS_MS` - время выполнения деления (мс).

Пример настройки (Windows):
```bash
set TIME_ADDITION_MS=1000
set TIME_MULTIPLICATIONS_MS=2000
set TIME_SUBTRACTION_MS=1000
set TIME_DIVISIONS_MS=3000
```

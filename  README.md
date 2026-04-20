## Структура проекта
**Стек:** Go 1.22 · `net/http` · `testify` · `allure-go` · `golangci-lint`

- **task1/** — результаты задания 1
    - `BUGS.md` — найденные баги с приоритизацией
- **task2/** — автотесты API
    - `internal/api/` — API-методы (создание, получение, удаление объявлений, статистика)
    - `internal/client/` — HTTP-клиент для работы с API
    - `internal/config/` — конфигурация проекта
    - `internal/models/` — структуры данных
    - `internal/suiteRun/` — базовый сьют для запуска тестов
    - `helpers/` — вспомогательные утилиты (assert)
    - `test/testdata/` — билдеры тестовых данных
    - `test/` — тестовые сценарии
    - `allure-report/` — сгенерированный отчёт Allure
    - `/task2/test/allure-results/` — результаты тестов для Allure
    - `BUGS.md` — найденные дефекты API
    - `TESTCASES.md` — описание тест-кейсов
    - `.golangci.yml` — конфигурация линтера

## Инструкция по установке и запуску



### Установка зависимостей

1. **Склонируйте репозиторий:**
   ```bash
   git clone https://github.com/zamurabims/QA_avito.git
   cd task2

2. **Установите зависимости:**
   ```bash
   go mod tidy

3. **Запуск тестов:**
   ```bash
   go test ./... -v

### С Allure-отчётом
1. **Установить Allure CLI: `brew install allure` или [allurereport.org](https://allurereport.org/)**

   ```bash
    ALLURE_RESULTS_DIR=./task2/test/allure-results go test ./task2/test/... -v
    allure generate ./task2/test/allure-results --clean -o ./task2/allure-report
    allure open ./task2/allure-report
   ```
---

## Линтер и форматтер

### Инструменты

| Инструмент       | Назначение                                      |
|------------------|-------------------------------------------------|
| `gofmt`          | Встроенный форматтер Go — выравнивание, отступы |
| `goimports`      | Форматтер + автоматическая сортировка импортов  |
| `golangci-lint`  | Агрегатор линтеров — запускает все проверки     |

### Установка

```bash
# goimports
go install golang.org/x/tools/cmd/goimports@latest

# golangci-lint (официальный способ)
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
  | sh -s -- -b $(go env GOPATH)/bin v1.57.0
```

### Команды

```bash
# Форматировать весь код (gofmt)
gofmt -w .

# Форматировать + упорядочить импорты (goimports)
goimports -w .

# Запустить все линтеры
golangci-lint run ./...

# Запустить с выводом в виде строк (удобно в CI)
golangci-lint run --out-format line-number ./...

# Проверить только конкретный пакет
golangci-lint run ./tests/create_item/...

# Автоматически исправить то что можно
golangci-lint run --fix ./...
```
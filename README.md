
# loglint

`loglint` — это анализатор для Go, который проверяет соблюдение правил написания сообщений логов.
Реализована поддержка "log/slog" и "go.uber.org/zap".

На данный момент проверяет:

* Сообщения логов должны начинаться со строчной буквы.
* Сообщения не должны содержать чувствительные данные (например, "password", "token").
* Сообщения должны быть на английском языке.
* Не допускаются специальные символы и эмодзи.

Линтер написан с использованием фреймворка `golang.org/x/tools/go/analysis`.

---

## Установка

### Сборка через go build

Клонирование:
```bash
git clone https://github.com/pollumna/loglint.git
cd loglint/cmd/loglint
```

Сборка исполняемого файла:
```bash
go build -o loglint
```

После этого в текущей директории появится файл loglint, который можно запускать:
```bash
./loglint ./...
```

### Через go install

#### Вариант 1. Клонирование репозитория
```bash
git clone https://github.com/pollumna/loglint.git
cd loglint/cmd/loglint
go install
```

#### Вариант 2. Без клонирования репозитория

```bash
go install github.com/pollumna/loglint/cmd/loglint@latest
```

После этого в `$GOBIN` появится исполняемый файл `loglint` и его можно запускать из любой директории:

```bash
loglint ./...
```

Например, на тестовых данных из директории проекта: 
```bash
loglint ./testdata/src/invalid/nonenglish/
```

---

## Использование

### На всём проекте

```bash
loglint ./...
```

### На конкретном файле

```bash
loglint path/to/file.go
```

### Пример неправильного и правильного кода

**Неправильно:**

```go
logger.Info("Starting server on port 8080")
slog.Error("ошибка")
```

**Правильно:**

```go
logger.Info("starting server on port 8080")
slog.Error("error")
```

---

## Тестирование

Тесты написаны с использованием `analysistest` из `golang.org/x/tools/go/analysis/analysistest`.

```bash
go test ./...
```

---

## Интеграция с golangci-lint

### Через плагин

Вы можете использовать `loglint` как плагин для golangci-lint. 
Для этого в корне проекта есть файлы `loglint.go`, `.custom-gcl.yml` и `.golangci.yaml`.

Для сборки необходимо запустить `golangci-lint custom` в корне проекта.
Полученный файл `./custom-gcl` - линтер golangci-lint с дополнительным плагином loglint.

---

### Через форк golangci-lint

Форк  golangci-lint с уже добавленным loglint https://github.com/pollumna/golangci-lint

Для сборки необходимо клонировать проект и выполнить:

```bash
go build -o <name>.exe ./cmd/golangci-lint
```

Полученный файл можно использовать в других проектах. Для корректного запуска в этом проекте собранного линтера 
необходимо удалить файлы `.custom-gcl.yml` и `.golangci.yaml`.

Тестирование:
```bash
.\<name>.exe run .\pkg\golinters\loglint\testdata\ --enable=loglint --default=none --max-same-issues=0
```

---

## Структура проекта

```
loglint/
 ├── analyzer/       # Реализация анализатора
 ├── cmd/loglint/    # CLI
 ├── testdata/       # Тестовые файлы для analysistest
 ├── loglint.go      # Плагин для golangci-lint
 ├── go.mod
 └── go.sum
```

---

## Планы

* Поддержка автофикса сообщений
* Расширяемые правила конфигурации

# log (`github.com/silverbulletsolutions/log`)

Lekka biblioteka do strukturalnego logowania zdarzeń w języku Go, wspierająca formaty JSON i konsolowy, interpolację parametrów w szablonach wiadomości, kontekst (`context.Context`) oraz zakresy pól (`Scope`).

## Instalacja

```bash
go get github.com/silverbulletsolutions/log
```

## Szybki start

### 1. Inicjalizacja loggera i podstawowe logowanie

```go
package main

import (
	"errors"

	"github.com/silverbulletsolutions/log"
)

func main() {
	// Dostępne implementacje: NewJsonLogger lub NewConsoleLogger
	// Poziomy: LogLevelTrace, LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError
	logger := log.NewJsonLogger(log.LogLevelInfo)

	// Logowanie z interpolacją parametrów {klucz}
	logger.Info("Użytkownik {user} wykonał akcję", log.Data{
		"user":    log.Str("jan_kowalski"),
		"user_id": log.Int(1234),
		"active":  log.Bool(true),
	})

	// Logowanie błędu
	err := errors.New("błąd połączenia")
	logger.Error("Wystąpił problem: {error}", log.Data{
		"error": log.Err(err),
		"retry": log.Int(3),
	})
}
```

### 2. Zakresy danych (`Scope`) i obsługa `context.Context`

Logger umożliwia tworzenie pod-loggerów z predefiniowanymi atrybutami (np. ID żądania) oraz przekazywanie ich przez standardowy `context.Context`:

```go
package main

import (
	"context"

	"github.com/silverbulletsolutions/log"
)

func processTask(ctx context.Context) {
	// Pobranie loggera z kontekstu
	if logger := log.FromContext(ctx); logger != nil {
		logger.Info("Przetwarzanie zadania", log.Data{
			"status": log.Str("w toku"),
		})
	}
}

func main() {
	baseLogger := log.NewConsoleLogger(log.LogLevelDebug)

	// Dodanie stałych pól do kontekstu loggera
	scopedLogger := baseLogger.Scope(log.Data{
		"request_id": log.Str("req-98765"),
	})

	// Osadzenie loggera w kontekście
	ctx := log.WithContext(context.Background(), scopedLogger)

	processTask(ctx)
}
```

### 3. Tworzenie loggera z konfiguracji (`CreateLogger`)

Funkcja `CreateLogger` pozwala zainicjalizować logger na podstawie obiektu implementującego interfejs `LogConfiguration` (odczytującego klucze `sw_logger` oraz `sw_logger_min_level`):

```go
package main

import "github.com/silverbulletsolutions/log"

type Config struct{}

func (c *Config) GetString(key string, defaultValue string) string {
	switch key {
	case "sw_logger":
		return "json" // "json" lub "console"
	case "sw_logger_min_level":
		return "info" // "trace", "debug", "info", "warn", "error"
	default:
		return defaultValue
	}
}

func main() {
	cfg := &Config{}
	logger := log.CreateLogger(cfg)
	logger.Info("Aplikacja uruchomiona", nil)
}
```

## Dostępne typy wartości (`LogValue`)

- `log.Str(string)` – ciąg znaków
- `log.Int(int)` / `log.UInt(uint)` – liczby całkowite
- `log.Float64(float64)` – liczby zmiennoprzecinkowe
- `log.Bool(bool)` – wartości logiczne
- `log.Err(error)` – formatowanie błędu jako ciąg znaków

*Wskazówka:* Klucze w `Data` zaczynające się od znaku `_` (np. `_tmp`) służą wyłącznie do interpolacji w treści wiadomości i nie są dołączane do atrybutów wyjściowych wpisu logu.

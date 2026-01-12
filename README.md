# go-config

Библиотека для загрузки конфигурации из YAML файлов с поддержкой подстановки переменных окружения.

## Установка

```bash
go get github.com/vodyanoysh/go-config
```

## Использование

### Базовый пример

```go
package main

import (
    "log"
    goconfig "github.com/vodyanoysh/go-config"
)

type Config struct {
    Database struct {
        Host     string `yaml:"host"`
        Port     int    `yaml:"port"`
        User     string `yaml:"user"`
        Password string `yaml:"password"`
    } `yaml:"database"`
    Server struct {
        Address string `yaml:"address"`
        Port    int    `yaml:"port"`
    } `yaml:"server"`
}

func main() {
    var config Config
    
    err := goconfig.LoadConfig(&config, "config.yaml")
    if err != nil {
        log.Fatal(err)
    }
    
    // Используйте config
    log.Printf("Database host: %s", config.Database.Host)
}
```

### Формат конфигурационного файла

Файл `config.yaml`:

```yaml
database:
  host: localhost
  port: 5432
  user: [[ getenv "DB_USER" ]]
  password: [[ getenv "DB_PASSWORD" ]]
server:
  address: [[ getenv "SERVER_ADDRESS" ]]
  port: 8080
```

### Загрузка .env файла

Библиотека автоматически загружает `.env` файл, если он существует в текущей директории:

```env
DB_USER=myuser
DB_PASSWORD=mypassword
SERVER_ADDRESS=0.0.0.0
```

## Формат подстановки переменных

Используйте формат `[[ getenv "ENV_NAME" ]]` для подстановки переменных окружения:

- `[[ getenv "DB_USER" ]]` - подставляет значение переменной окружения `DB_USER`
- Можно использовать в любом месте строки: `prefix [[ getenv "VAR" ]] suffix`
- Поддерживается множественная подстановка в одной строке

Если переменная окружения не найдена, будет выведено предупреждение в лог, а значение останется пустым.


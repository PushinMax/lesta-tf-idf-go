# TF-IDF Document Analyzer

Сервис для анализа текстовых документов с использованием алгоритма TF-IDF. Позволяет загружать документы, анализировать частотность слов и управлять коллекциями документов.

## Features

- 📊 Анализ документов с помощью TF-IDF
- 🔐 JWT аутентификация
- 📁 Управление коллекциями документов
- 📈 Статистика использования слов
- 🗜️ Сжатие текста алгоритмом Хаффмана
- 📊 Метрики использования API

## Tech Stack

- Go (Gin Framework)
- PostgreSQL
- MongoDB
- nGinx
- Docker & Docker Compose
- JWT Authentication
- Swagger Documentation

## Quick Start

### Предварительные требования

- Docker и Docker Compose
- Нормальная операционная система(не windows)

### Установка и запуск

1. Клонируйте репозиторий:
```bash
git clone https://github.com/PushinMax/lesta-tf-idf-go.git
cd lesta-tf-idf-go
```

2. Сгенерируйте файл конфигурации:
```bash
sh scripts/generate_env.sh
```

3. Запустите приложение:
```bash
docker-compose up -d --build
```

## API Documentation

API документация доступна по адресу `/swagger/index.html` после запуска приложения.


## Демо

Демо версия приложения доступна по адресу: [37.9.53.117:80](http://37.9.53.117:80)

## Development

### Структура проекта

```
.
├── cmd/            # Точка входа приложения
├── internal/       # Внутренняя логика
│   ├── encoding/   # Логика сжатия документов
│   ├── handler/    # HTTP обработчики
│   ├── repository/ # Слой работы с БД
│   ├── server/     # Запуск сервера
│   └── service/    # Бизнес-логика
├── migrations/     # SQL миграции
└── scripts/        # Вспомогательные скрипты
```

## License

MIT

## Changelog

С changelog-ом вы можете ознакомиться [тут](docs/CHANGELOG.md) и дополнить его своим комментарием

## Schema Database

### PostgreSQL Users Table
```mermaid
erDiagram
    USERS {
        uuid id PK
        varchar(20) login
        text password_hash
        text token_hash
        timestamptz created_at
    }
```

### MongoDB Collections

#### Document Structure
```mermaid
graph LR
    Document{Document} --> ID[ObjectID]
    Document --> FileInfo[File Info]
    Document --> Content[Content]
    Document --> Statistics[Statistics]
    Document --> Collections[Collections]
    
    FileInfo --> file_id[file_id: string]
    FileInfo --> file_name[file_name: string]
    FileInfo --> user_id[user_id: string]
    
    Content --> file[file: string]
    Content --> len[len: int]
    
    Statistics --> stats[stats: WordStat[]]
    Statistics --> words[words: Map<string,int>]
    Statistics --> isvalid[isvalid: bool]
    
    Collections --> collections_list[collections: string[]]
    Document --> huffman[huffman_encoding: HuffmanCode]
```

#### Collection Structure
```mermaid
graph LR
    Collection{Collection} --> ID[ObjectID]
    Collection --> Info[Info]
    Collection --> Statistics[Statistics]
    Collection --> Documents[Documents]
    
    Info --> name[name: string]
    Info --> user_id[user_id: string]
    
    Statistics --> stats[stats: WordStat[]]
    Statistics --> words[words: Map<string,WordCount>]
    Statistics --> isvalid[isvalid: bool]
    Statistics --> len[len: int]
    
    Documents --> documents_id[documents: string[]]

    WordCount{WordCount} --> amount_w[amount_w: int]
    WordCount --> amount_d[amount_d: int]
```

## Contact

При возникновении вопросов или проблем создавайте issue в репозитории проекта.
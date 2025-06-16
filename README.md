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
        varchar login
        text password_hash
        text token_hash
        timestamp created_at
    }
```

### MongoDB Collections

#### Document Structure
```mermaid
flowchart LR
    Document[Document] --> ID[_id: ObjectID]
    Document --> FileInfo[File Information]
    Document --> Content[Content]
    Document --> Stats[Statistics]
    Document --> CollectionsList[Collections List]
    
    FileInfo --> FID[file_id: string]
    FileInfo --> FName[name: string]
    FileInfo --> UID[user_id: string]
    
    Content --> Text[content: string]
    Content --> Length[length: int]
    
    Stats --> WordStats[stats: Array]
    Stats --> WordCount[words: Map]
    Stats --> Valid[isValid: boolean]
    
    CollectionsList --> CList[collections: Array]
    Document --> Huffman[huffman: Object]
```

#### Collection Structure
```mermaid
flowchart LR
    Collection[Collection] --> CID[_id: ObjectID]
    Collection --> BasicInfo[Basic Info]
    Collection --> Stats[Statistics]
    Collection --> Docs[Documents]
    
    BasicInfo --> Name[name: string]
    BasicInfo --> UID[user_id: string]
    
    Stats --> WordStats[stats: Array]
    Stats --> WordMap[words: Map]
    Stats --> Valid[isValid: boolean]
    Stats --> Length[length: int]
    
    Docs --> DocList[documents: Array]
    
    WordMap --> Counter[WordCounter]
    Counter --> DocCount[doc_count: int]
    Counter --> WordCount[word_count: int]
```

## Contact

При возникновении вопросов или проблем создавайте issue в репозитории проекта.
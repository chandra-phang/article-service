# Article Service

> A Go-based backend service for creating and listing articles, with PostgreSQL for primary data storage and Elasticsearch for full-text search.

## Features

- Create and list articles
- Search articles by title, body, or author name using PostgreSQL and Elasticsearch
- Integration and unit testing with test database and factory data

## Tech Stack

- **Language:** Go
- **Database:** PostgreSQL
- **Search Engine:** Elasticsearch 8.12.0
- **Testing:** `sqlmock`, Go `testing` package
- **Logging:** `logrus` for structured logging and `lumberjack` for log file rotation
- **UUID:** `github.com/google/uuid`
- **Elasticsearch Client:** `github.com/olivere/elastic/v7`

## Getting Started

### 1. Clone the repo

```bash
git clone https://github.com/yourusername/article-service.git
cd article-service
```

### 2. Setup PostgreSQL

```bash
CREATE DATABASE article_service;
```

### 3. Setup Elasticsearch

```bash
docker run -d --name elasticsearch -p 9200:9200 -e "discovery.type=single-node" elasticsearch:8.12.0
```

### 4. Configuration

Update your database and Elasticsearch connection in `config.yml`:

```yaml
app:
  port: ":3000"

db:
  host: "localhost"
  port: 5432
  user: your_user
  password: your_password
  dbname: "article_service"
  test_dbname: "article_service_test"

elastic:
  url: "http://localhost:9200"
```

### 5. Run the application

```bash
go run main.go
```

## API Endpoints

| Method | Endpoint      | Description             |
| ------ | ------------- | ----------------------- |
| POST   | `v1/articles` | Create a new article    |
| GET    | `v1/articles` | List or search articles |

## Running Tests with Makefile

The project includes a `Makefile` for running tests and generating code coverage reports across platforms.

### Run Tests (Windows)

```bash
make test
```

### Run Tests (Linux/macOS)

```bash
make test-unix
```

### View HTML Coverage Report

```bash
make coverage
```

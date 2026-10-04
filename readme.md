# URL Shortener

A simple URL shortening service built with **Go**, **PostgreSQL**, and **Redis**.

## Features

* Create short URLs
* Redirect short URLs to original URLs
* Redis caching for faster lookups
* PostgreSQL for persistent storage
* REST API
* Middleware and graceful shutdown

## Tech Stack

* **Go** — Backend
* **PostgreSQL** — Database
* **Redis** — Caching
* **Docker Compose** — Redis deployment

## API

### Create Short URL

```http
POST /shorten
```

Request:

```json
{
  "url": "https://example.com"
}
```

### Redirect

```http
GET /abc123
```

Redirects to the original URL.

## Run Locally

Clone the repository:

```bash
git clone https://github.com/shouvik214/Url_shortener_project.git
cd Url_shortener_project
```

Start Redis:

```bash
docker compose up -d
```

Then run the Go server:

```bash
go run ./cmd/server
```

The server runs on:

```text
http://localhost:8080
```

## Author

**Shouvik Mondal**

GitHub: https://github.com/shouvik214

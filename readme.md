# URL Shortener

A simple and scalable URL shortening service built with **Go**, **PostgreSQL**, **Redis**, and **Docker**.

## Features

* Create short URLs
* Redirect short URLs to original URLs
* Redis caching for faster lookups
* PostgreSQL for persistent storage
* REST API
* Docker support
* Graceful shutdown and middleware

## Tech Stack

* **Go** — Backend
* **PostgreSQL** — Database
* **Redis** — Caching
* **Docker** — Containerization

## API

### Create Short URL

POST /shorten

Request:

{
  "url": "https://example.com"
}


Response:

{
  "short_url": "http://localhost:8080/abc123"
}

### Redirect


GET /abc123

Redirects to the original URL.

## Run Locally

Clone the repository:

git clone https://github.com/shouvik214/Url_shortener_project.git
cd Url_shortener_project

Run with Docker:

docker compose up --build

run directly with Go:

go run ./cmd/server

The server runs on:
http://localhost:8080

## Author

**Shouvik Mondal**

[GitHub](https://github.com/shouvik214)

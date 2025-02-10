# Hack the North 2025 Backend Challenge

This is my backend API for the take-home assignment! For my tech stack, I used Go, SQLite w/ GORM, and Chi. I decided to use SQLite because it's a lightweight database that's easy to set up and use (and it was also recommended), and GORM because it made database interaction simple. Furthermore, I used Chi since I like the ability to mount multiple routers on a single server, which allows me to organize controllers more easily.

For assumptions, I didn't assume too much other than the clarifications given to us (ex. times given in ISO 8601, emails & badge codes are unique, etc.). I otherwise try account for anything else that was unclear. For example, I account for empty badge codes by keeping it as a nullable field in the database.

In terms of database structure, I have two models: `User` and `Scan`. The `User` model contains fields for the user's name, email, phone number, badge code, and scans. The `Scan` model contains fields for the scan's activity name, category, and time. These two models are linked by a one-to-many relationship, where a user can have many scans.

As specified in the challenge, the backend also keeps track of the last updated time for each user. This is mainly handled by GORM since it already keeps track of such a field, but it is also manually updated in the `CreateScan` function.

## Setup

1. Clone the repository

```bash
git clone https://github.com/AritroSaha10/htn25-backend-takehome.git
```

2. Install the dependencies

```bash
go mod tidy
```

3. Run the server

```bash
go run .
```

## API Documentation

The API documentation is available at `http://localhost:8080/swagger/index.html`.

## Docker

1. Build the Docker image

```bash
docker build -t htn25-backend-takehome .
```

2. Run the Docker container using Compose

```bash
docker compose -f docker-compose.dev.yml up
```

## Author

Written by Aritro Saha.

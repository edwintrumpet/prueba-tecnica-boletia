# Boletia currencies

App to know the historical behavior of currencies.

## Develop

Requirements

- Go v1.20+
- Docker
- Docker compose

Create the `.env` file to provide env variables for database

```
POSTGRES_PASSWORD=******
POSTGRES_USER=******
POSTGRES_DB=******
```

Run database using docker compose

```bash
docker-compose up -d
```

Create the `.vscode/launch.json` file to provide debug config and env variables

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Boletia currencies",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd",
            "env": {
                "REQUESTS_TIME": "480",
                "TIMEOUT": "30",
                "CURRENCIES_HOST": "https://api.currencyapi.com/v3/latest",
                "API_KEY": "******",
                "DB_USER": "******",
                "DB_PASSWORD": "******",
                "DB_NAME": "******",
            }
        }
    ]
}
```

To run press **`F5`** or use the debug icon on Visual Studio Code

### Run using Docker

Add next lines to the `.env` file to provide env variables

```
REQUESTS_TIME=480
TIMEOUT=30
CURRENCIES_HOST=https://api.currencyapi.com/v3/latest
API_KEY=******
DB_USER=******
DB_PASSWORD=******
DB_NAME=******
```

Enable commented lines on `docker-compose.yaml` file

Run

```bash
docker-compose up -d --build app
```

Watch app logs

```bash
docker-compose logs -f app
```

## Tests

Run tests

```bash
docker-compose -f docker-compose.test.yaml up -d && sleep 5 && go test -v ./... ; docker-compose down
```

Watch coverage

```bash
docker-compose -f docker-compose.test.yaml up -d && sleep 5 && go test ./... -coverprofile cover.out && go tool cover -func cover.out ; docker-compose down
```

Watch coverage in html

```bash
docker-compose -f docker-compose.test.yaml up -d && sleep 5 && go test ./... -coverprofile cover.out && go tool cover -html=cover.out ; docker-compose down
```

## Deploy

Create a git tag using semantic version

```bash
git tag -a v0.0.1 -m "message"
```

Push tag to GitHub

```bash
git push --tags
```

To see all versions visit [Docker Hub repository](https://hub.docker.com/repository/docker/edwincoding/boletia-currencies/)

The server must have the next `docker-file.yaml`

```yaml
services:
  app:
    image: edwincoding/boletia-currencies:$VERSION
    environment:
      - VERSION
      - REQUESTS_TIME
      - TIMEOUT
      - CURRENCIES_HOST
      - API_KEY
      - DB_USER
      - DB_PASSWORD
      - DB_NAME
      - DB_HOST=db
    depends_on:
      db:
        condition: service_healthy
    ports:
      - 80:3000
  db:
    image: postgres
    environment:
      - POSTGRES_PASSWORD
      - POSTGRES_USER
      - POSTGRES_DB
    healthcheck:
      test: ["CMD-SHELL", "pg_isready"]
      interval: 10s
      timeout: 5s
      retries: 5
    volumes:
      - ./pgdata:/var/lib/postgresql/data
```

## To remove a version

Delete tag on [Docker Hub repository](https://hub.docker.com/repository/docker/edwincoding/boletia-currencies/)

Remove tag from GitHub

```bash
git push --delete origin v0.0.1
```

Remove tag from local repository

```bash
git tag --delete v0.0.1
```

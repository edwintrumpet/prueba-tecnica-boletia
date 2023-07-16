# Boletia currencies

App to know the historical behavior of currencies.

## Develop

Requirements

- Go v1.20+

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
            }
        }
    ]
}
```

To run press **`F5`** or use the debug icon on Visual Studio Code

### Run using Docker

Create the `.env` file to provide env variables

```
REQUESTS_TIME=480
TIMEOUT=30
CURRENCIES_HOST=https://api.currencyapi.com/v3/latest
API_KEY=******
```

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
go test -v ./...
```

Watch coverage

```bash
go test ./... -coverprofile cover.out && go tool cover -func cover.out
```

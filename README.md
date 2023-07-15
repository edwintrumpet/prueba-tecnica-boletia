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
                "REQUESTS_TIME": "1",
                "TIMEOUT": "60",
                "CURRENCIES_HOST": "http://localhost:3000",
            }
        }
    ]
}
```

To run press **`F5`** or use the debug icon on Visual Studio Code

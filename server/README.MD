# Websocket server

## Getting project up
```bash
docker-compose build && docker-compose up
```

## Configuration
Change environmental variables inside docker-compose or create docker-composer.override.yml and change there


## Decision explanations

Most of the wiring would go to routes.go ( unless project would be bigger and then it would need some dependency injection framework)

Logger is initiated inside main, as server depends on it.


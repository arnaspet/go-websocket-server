# Websocket server

## Getting project up
```bash
docker-compose up
```

you can also scale subscribers or publishers like so:
i want 1 publishers and 5 subscribers

```bash
docker-compose up --scale subscriber=5
```

## TODO
* Publish only when there is subscribers listening on server
* Publisher & subscriber better logging
* Fix race condition on exiting websocket
* Reuse logic of publisher & subscriber with composition
* Add tests
# Teso task

Using websockets publisher sends a message every second to server if there is subscribers.
Server fanouts them to subscribers.

## Flow:

```
                                                                                                          +------------------------+
                                                                                                          |                        |
                                                                                                +----->   |       Subscriber       |
                                                                                                |         |                        |
                                                                                                |         +------------------------+
                                                                                                |
+------------------------+                     +------------------------+                       |         +------------------------+
|                        |   Msg: Yellow?      |                        |         Msg: Yellow!  |         |                        |
|        Publisher       | +---------------->  |         Server         |  +-------------------------->   |       Subscriber       |
|                        |                     |                        |                       |         |                        |
+------------------------+                     +------------------------+                       |         +------------------------+
                                                                                                |
                                                                                                |         +------------------------+
                                                                                                |         |                        |
                                                                                                +----->   |       Subscriber       |
                                                                                                          |                        |
                                                                                                          +------------------------+

```

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
* Add context logger to server
* Add tests

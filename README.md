## .env Example

```dotenv
LOCATION=remote
####################################################
###################### LOGGER ######################
LOGGER_LEVEL=debug
LOGGER_ENCODING=console
# zapcore.LevelEncoder.UnmarshalText : capital | capitalColor | color | lower
LOGGER_LEVEL_ENCODER=lower
####################################################
##################### DATABASE #####################
# LOCAL
DB_DSN=postgres://todo:todo@localhost:5490/todo?sslmode=disable
```

## Run application

```bash
$ docker-compose -f ./ops/docker/docker-compose.yaml
```

And if you see the following output then services is started ok.
```bash
Running 3/3
  Container docker-todo      Started                                                                                                                6.3s
```

Run todo service
```bash
go run ./cmd
```

## Swagger API

swagger-ui looks like this:
![Demo-Api](swagger-ui.png)
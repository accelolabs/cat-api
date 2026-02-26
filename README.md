## Что это? 

Сервис замяукивания ссылок. Превращает любую ссылку в мяукания.

[meow.accelolabs.com](https://meow.accelolabs.com/)

Пример ссылки - [https://meow.accelolabs.com/mrrrp-meeeooww-mrroooowww-purrrrr-maaaowwww-mrroow](https://meow.accelolabs.com/mrrrp-meeeooww-mrroooowww-purrrrr-maaaowwww-mrroow)

## Что умеет?

Умеет замяукивать ссылки. Псевдонимы для ссылок хранятся в SQLite базе данных.

Есть минималистичный фронтенд для этого сервиса ([web/index.html](https://github.com/accelolabs/cat-api/blob/main/web/index.html)). Он может быть указан в конфиге. Тогда будет использоваться root эндпойнт для сервировки html страницы.

## Как использовать?

Сбилдить проект можно командой:
```
go build -o meow ./cmd/cat-api/main.go
```

При запуске конфиг нужно указать аргументом, пример конфига лежит в [config/local.yaml](https://github.com/accelolabs/cat-api/blob/main/config/local.yaml)
```
./meow -config=./config/local.yaml
```

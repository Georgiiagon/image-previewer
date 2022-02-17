# Превьювер изображений

Сервис предназначен для изготовления preview (создания изображения с новыми размерами на основе имеющегося изображения).

## Установка:
1. Создание .env файла
```shell
cp env.example .env
```
2. Запуск приложения через докер
```shell
make run-img
```

## API

Отправьте GET запрос по следующему адресу
```shell
http://localhost:4000/resize/{width}/{height}/{http://image.jpg}
```
Например
```shell
http://localhost:3000/resize/100/400/http://sebweo.com/wp-content/uploads/2020/01/what-is-jpeg_thumb-800x478.jpg?x72922
```

Доступные make команды:
1. `make build`
2. `make run`
3. `make build-img`
4. `make run-img`
5. `make test`
6. `make lint`


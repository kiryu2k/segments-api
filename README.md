# Segments API
Segments API — сервис, хранящий пользователей и сегменты, в которых они состоят (создание, удаление сегментов, а также добавление и удаление пользователей в сегмент).
## Usage
Для поднятия и развертывания dev-среды выполнить:
```
docker-compose -f docker-compose.dev.yaml up
```
Локально:
```
docker-compose -f docker-compose.local.yaml up
```
Сборка и запуск приложения (по умолчанию загружается конфигурация для dev-среды):
```
make
```
Сборка и запуск приложения с конфигурацией для local-среды:
```
make build
./bin/segments --config ./configs/config.local.yaml
```
Изменить конфигурацию для той или иной среды можно в файлах `config.dev.yaml` и `config.local.yaml` в директории `./configs`. Также обязательно создать `.env` файл с необходимыми переменными окружения (смотри `example.env`).
## Endpoints
**Swagger документация**:
```
GET /docs/index.html
```
**Метод создания сегмента.** Принимает в body slug (название) сегмента. Опционально можно указать процент пользователей, которые добавятся в этот сегмент автоматически:
```
POST /segment
```
**Метод удаления сегмента.** Принимает slug (название) сегмента:
```
DELETE /segment/{slug}
```
**Метод создания пользователя.** Принимает в body id пользователя:
```
POST /user
```
**Метод изменения активных сегментов пользователя.** Принимает в body список slug (названий) сегментов которые нужно добавить пользователю, 
список slug (названий) сегментов которые нужно удалить у пользователя, id пользователя. Также есть возможность задать TTL для добавляемых сегментов, 
чтобы по истечению времени они автоматически удалились у пользователя. TTL задается в формате "1y8m21d". 
Если хотите только удалить определенные сегменты, то можно опустить список сегментов для добавления и наоборот:
```
POST /user-segments
```
**Метод получения активных сегментов пользователя.** Принимает на вход id пользователя:
```
GET /user-segments/{userID}
```
**Метод удаления пользователя.** Принимает на вход id пользователя:
```
DELETE /user/{userID}
```
**Метод получения истории добавления и удаления** сегментов указанного пользователя за определенные год и месяц (указываются в query параметрах в численном виде) в формате CSV:
```
GET /log/{userID}?date={year-month}
```
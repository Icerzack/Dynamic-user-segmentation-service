[![codecov](https://codecov.io/gh/Icerzack/Dynamic-user-segmentation-service/branch/main/graph/badge.svg)](https://codecov.io/gh/Icerzack/Dynamic-user-segmentation-service)
# Dynamic user segmentation service
Сервис динамического сегментирования пользователей
<details>
  <summary>Содержание</summary>
  <ol>
    <li><a href="#основные-сведения">Основные сведения</a></li>
    <li><a href="#установка-и-запуск">Установка и запуск</a></li>
    <li><a href="#тесты-и-ci">Тесты и CI</a></li>
  </ol>
</details>

## Основные сведения
### Описание
Данная программа содержит реализацию сервиса динамического сегментирования пользователей. 

Ниже приведены основные моменты данной реализации:
 1. Реализовано на **Golang**;
 2. Есть возможность экспорта **истории добавления\удаления** сегментов у пользователей;
 3. Есть поддержка назначения **TTL** определенному сегменту у пользователя;
 4. Используются следующие библиотеки:

    2.1. https://github.com/georgysavva/scany -- Работа с БД;
    
    2.2. https://github.com/jackc/pgconn -- Работа с БД;
    
    2.3. https://github.com/jackc/pgx/ -- Работа с БД;
    
    2.4. https://pkg.go.dev/go.uber.org/mock -- Создание stub'ов и mock'ов для тестов.
    
 5. **docker-compose** для поднятия и развертывания среды;
 6. **CI** на **GitHub Actions** и тесты для проверки главного функционала.

Структура проекта выглядит следующим образом:
```markdown
├── .github
|   └── workflows
|       └── tests.yml                   - Файл конфигурации CI
├── build
|   └── Dockerfile                      - Dockerfile для создания образа приложения
├── docs
|   └── history.csv                     - CSV файл, содержащий в себе историю по операциям
├── cmd
|   └── main.go                         - Запускающий файл
├── internal
|   ├── app
|   |   ├── config.go                   - Конфигурация сервера
|   |   ├── server.go                   - Методы для запуска сервера и обработки запросов
|   |   └── server_test.go              - Тесты для проверки работоспособности
|   └── pkg
|       ├── db
|       |   ├── mocks
|       |   |   └── mock.go             - Авто-генерация моков для интерфейса DBops
|       |   ├── sql
|       |   |   └── init.sql            - Файл для создания таблиц в базе данных
|       |   ├── client.go               - Создание экземпляра PostgreSQL
|       |   └── database.go             - Методы по работе с базой данных
|       ├── model
|       |   ├── segment.go              - Структуры для работы с JSON
|       |   └── user_segment.go         - Структуры для работы с JSON
|       ├── repository
|       |   ├── postgresql
|       |   |   ├── segments.go         - Реализация паттерна репозиторий для PostgreSQL
|       |   |   └── user_segment.go     - Реализация паттерна репозиторий для PostgreSQL
|       |   ├── repository.go           - Паттерн репозиторий
|       |   └── structs.go              - Структуры для базы данных
|       └── service
|           ├── db
|           |   ├── mocks
|           |   |   └── mock.go         - Авто-генерация моков для интерфейса Service
|           |   ├── postgres.go         - Реализация бизнес-логики для PostgreSQL
|           |   └── service.go          - Интерфес для работы бизнес-логики
|           └── history.go
|               ├── csv.go              - Реализация модуля ведения истории для формата CSV
|               └── service.go          - Интерфес для работы c модулем истории
├── .gitignore
├── docker-compose.yml                  - Для поднятия в Docker'е
├── go.mod
├── Makefile
└── README.md
```
### База данных
**Подразумевается**, что ID пользователя нам известен заранее и подается из вне как запрос на этот сервис.

![image](https://github.com/Icerzack/avito-backend-internship/assets/24461208/94573a6a-676d-43d9-b003-a28d5e68b5ef)

Таблица **segments** хранит в себе в качестве ключа название сегмента и его описание. Описание может быть пустым.

Таблица **user_segments** хранит в себе в качестве ключа **id** записи в этой таблице, **user_id** хранит в себе уникальный идентификатор пользователя (который мы знаем заранее), **seg_title** хранит в себе название сегмента, в котором состоит пользователь. **seg_title** не является **FK**! Данная зависимость реализована напрямую из кода, дабы избежать ссылок в самих таблицах.   

Таким образом, чтение, добавление и удаление сегментов у пользователя реализовано по следующему алгоритму:
1. Если данный ID пользователя имеется в таблице **user_segments**, то будут выведены все сегменты из данной таблицы, где поле **user_id** равняется нашему запрашиваемому;
2. Если данный ID пользователя в таблице НЕ имеется, то соответственно количество найденных строк будет 0. Поэтому, в качестве ответа, вернется сообщение о том, что данный пользователь не состоит ни в одном сегменте;
3. Если при добавлении пользователю перечислены несуществующие в таблице **segments** названия сегментов, то записи не будут созданы в таблице **user_segments**, а в качестве ответа вернется массив элементов, которые не были добавлены из-за этого несоответсвия. В остальном в таблицу будут добавлены новые записи, содержащие в себе ID пользователя и сегмент (и так для каждого из перечисленных существующих сегментов);
4. Если при удалении у пользователя перечислены несуществующие в таблице **user_segments** названия сегментов, то подходящие записи не найдутся, а следовательно, в качестве ответа, будет возвращен массив элементов, которые удалить не удалось. В остальном у данного пользователя будут убраны соотвествующие сегменты из данной таблицы.

### Ведение истории операций
При выполнении операции добавления\удаления в файле **(docs/history.csv)** автоматически будут появляться записи в формате CSV. Чтобы посмотреть данный файл, нужно перейти по адресу: **localhost:3001/docs/history.csv** (если запущено в Docker).

Пример вывода:

<img width="250" alt="Снимок экрана 2023-08-31 в 22 26 43" src="https://github.com/Icerzack/avito-backend-internship/assets/24461208/6d3e5e76-5ef5-4103-9118-089c1e84ad06">

### Взаимодействие
Для взаимодействия с сервером есть 4 способа:
1. /create-segment -- создать сегмент. Принимает JSON в качестве тела запроса;
2. /delete-segment -- удалить сегмент. Принимает JSON в качестве тела запроса;
3. /user-in-segment -- выполнить операции создания и\или удаления для определенного пользователя. Принимает JSON в качестве тела запроса;

   3.1. Поле **"ttl"** передается как массив **int**-ов, Так, элемент на позиции **i** будет назначен элементу на позиции **i** в массиве **"seg_titles_to_add"**. Если значение 0, тогда назначение TTL игнорируется. Передается в секундах.
5. /get-user-segments -- получить все сегменты, в которых состоит пользователь. Принимает JSON в качестве тела запроса.

В качестве ответа возвращается JSON, который содержит всегда поле "status" и дополнительные.

![image](https://github.com/Icerzack/avito-backend-internship/assets/24461208/decd04eb-4df5-4e69-b8a2-b13a96e9bf6b)

Сами структуры запросов и ответов в этих файлах:
1. https://github.com/Icerzack/avito-backend-internship/blob/main/internal/pkg/model/segment.go
2. https://github.com/Icerzack/avito-backend-internship/blob/main/internal/pkg/model/user_segment.go

### Примеры
Примеры запросов и ответов в Postman:
1. Создать сегмент (корректный запрос):

![image](https://github.com/Icerzack/avito-backend-internship/assets/24461208/ce770c24-ae46-4cce-bc5e-94651bb0737d)

2. Удалить сегмент (корректный запрос):

![image](https://github.com/Icerzack/avito-backend-internship/assets/24461208/22c96dad-b457-4976-9a1f-a1f3c117ebdf)

3. Создать сегмент (некорректный запрос, неправильное наименование полей):

![image](https://github.com/Icerzack/avito-backend-internship/assets/24461208/44d39f69-7b91-4166-8ecd-2b9a97fac627)

4. Добавить пользователя в сегмент (корректный запрос):

![image](https://github.com/Icerzack/avito-backend-internship/assets/24461208/32b7f640-34d9-481b-9a05-953cd1b72968)

5. Добавить пользователя в сегмент (корректный запрос, но некоторые названия не существуют):

![image](https://github.com/Icerzack/avito-backend-internship/assets/24461208/1ba32b96-9a6d-4631-9162-e5848f12854b)

6. Удалить пользователя из сегмента (корректный запрос):

![image](https://github.com/Icerzack/avito-backend-internship/assets/24461208/42b9b7a9-bfb8-4693-b9f4-79e13563b083)

7. Удалить пользователя из сегмента (корректный запрос, но некоторые названия не существуют):

![image](https://github.com/Icerzack/avito-backend-internship/assets/24461208/3f4691ee-882a-444c-9baf-be610cd1db47)

8. Одновременное создание и удаление (корректный запрос, но некоторые названия не существуют):

![image](https://github.com/Icerzack/avito-backend-internship/assets/24461208/e2a1b22a-e1ea-43f1-84b2-944205b2cc7b)

9. Создать сегменты с временем жизни (корректный запрос)

![image](https://github.com/Icerzack/avito-backend-internship/assets/24461208/43047fa1-9402-4dbd-b21e-1ee249c438b2)

10. Получить сегменты пользователя (корректный запрос):

![image](https://github.com/Icerzack/avito-backend-internship/assets/24461208/ca02f8b2-a54c-4155-a3ba-87c6a03b8647)

11. Получить сегменты пользователя (корректный запрос, но данный пользователь не содержит сегментов):

![image](https://github.com/Icerzack/avito-backend-internship/assets/24461208/803f39a2-b7b7-4c59-b1d7-024691185a63)

## Установка и запуск

Проект содержит **Makefile**, который имеет следующий вид:

```makefile
.PHONY: compose-db-up
compose-db-up:
	docker-compose build
	docker-compose up -d postgres

.PHONY: compose-db-rm
compose-db-rm:
	docker-compose down

.PHONY: compose-app-up
compose-app-up:
	docker build -f build/Dockerfile -t avito-app .
	docker-compose up -d avito-app

.PHONY: compose-app-rm
compose-app-rm:
	docker-compose down

.PHONY: compose-all-up
compose-all-up: compose-db-up compose-app-up

.PHONY: compose-all-rm
compose-all-rm: compose-app-rm compose-db-rm

.PHONY: run-tests
run-tests:
	go test -v avito-backend-internship/internal/app
```
Для запуска конкретно одного из компонентов приложения нужно, например, прописать следующее:
```bash
$ make compose-db-up
```
Для поднятие всего целиком:
```bash
$ make compose-all-up
```

## Тесты и CI
Для тестирования использовалась библиотека **gomock**. С ее помощью были созданы моковые реализации БД и Сервиса на основании их интерфейсов. Данные реализации применяются в файле **(internal/app/handlers_test.go)**, который тестирует главный функционал приложения (корректная обработка URL и добавление в БД).

В проекте был настроен простой CI для запуска кода на тестах.
Файл конфигурации может быть найден по следующему пути: **(.github/workflows/tests.yml)**
и имеет следующий вид:
```yml
name: tests

on:
  push:
    branches: [ "main" ]
  pull_request:
    branches: [ "main" ]

jobs:

  build:
    runs-on: ubuntu-latest
    steps:
    - name: Checkout repository
      uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.19

    - name: Tests
      run: go test -v avito-backend-internship/internal/app
 ```

<p align="right">(<a href="#основные-сведения">К началу</a>)</p>

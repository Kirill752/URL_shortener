<h3 align="center">URL SHORTENER</h3>

<div align="center">

[![Status](https://img.shields.io/badge/status-active-success.svg)]()

</div>

---

<p align="center"> 
URL SHORTENER - это небольшой REST API сервис для собственного пользования, созданный для того, чтобы упростить поиск web страниц, адрес которых имеет довольно сложную структуру.
    <br> 
</p>

## Оглавление

- [О проекте](#about)
- [Начало работы](#getting_started)
- [Запуск тестов](#tests)
- [Использовнаие](#usage)
- [Использованные библиотеки](#built_using)
- [Автор](#authors)
- [Благодарности](#acknowledgement)

## О проекте <a name = "about"></a>

URL SHORTENER(сокращатель ссылок) представляет из себя сервис для удобного и быстрого доступа к сайтам, чей адрес является сложным и нетривиальным. Например, требуется перейти на специфическую страницу по типу:

https://multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/

Но вместо этого можно один раз занести данный URL в базу данных локально развернутого сервиса и каждый раз при надобности переходить по ссылке вида:

http://localhost:8082/chip-8

## Начало работы <a name = "getting_started"></a>

### Предварительные действия

Перед началом работы с сервисом убедитесь, что у вас есть компилятор Golang.
Как его установить написано на официальном [сайте](https://go.dev/).

### Установка

Данный сервис использует сторонние библиотеки. Они указанны в разделе 
[Использованные библиотеки](#built_using)

Если какого то модуля нет, необходимо добавить его с помощью команды go get [link].

## Запуск тестов <a name = "tests"></a>

Для запуска тестов необходимо выполнить команду:
```
make test
```

### Что проверяют тесты?

Тесты сосотоят из тестов хэндлеров, а также тестов сервиса как черной коробки.

## Использовнаие <a name="usage"></a>

Для запуска сервера необходимо выполнить команду:
```
make run
```
Для отсановки работы сервера необходимо выполнить команду:
```
make kill
```
Для добаления скоращения для ссылки необходимо выполнить команду:
```
make add URL=[your-url] ALIAS=[your alias]
```
Для удаления записи из базы данных необходимо выполнить команду:
```
make delete ALIAS=[your alias]
```
Для перехода на сайт по существующим алиасам необходимо в открытом браузере
перейти по адресу: http://localhost:8082/{your_alias}

## Использованные библиотеки <a name = "built_using"></a>

- [SQLite](https://www.sqlite.org/) - База данных
- [Chi](https://go-chi.io/#/) - Обработка HTTP запросов
- [Testify](https://pkg.go.dev/github.com/stretchr/testify) - тестирование
- [Cleanenv](https://pkg.go.dev/github.com/ilyakaznacheev/cleanenv) - конфигурирование

## Благодарности <a name = "acknowledgement"></a>

- При создании данного сервиса мне очень сильно помогла [данная статья](https://habr.com/ru/companies/selectel/articles/747738/).

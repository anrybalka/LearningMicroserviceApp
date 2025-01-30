# Learning Microservice App

Этот проект представляет собой пример микросервисного приложения, состоящего из нескольких сервисов, взаимодействующих друг с другом. Каждый сервис выполняет свою уникальную задачу и работает в изолированной среде с использованием Docker.

## Содержание
1. [Структура проекта](#структура-проекта)
2. [Запуск проекта](#запуск-проекта)
3. [Описание сервисов](#описание-сервисов)
4. [API документация](#api-документация)
5. [Примеры запросов](#примеры-запросов)

---

## Структура проекта

Проект состоит из следующих сервисов:
- **nginx**: Веб-сервер для маршрутизации запросов.
- **000-user-service**: Сервис для регистрации, аутентификации и проверки токенов пользователей.
- **001-log-service**: Сервис для логирования событий.
- **002-robot-tester**: Сервис для автоматического тестирования функциональности других сервисов.


## Запуск проекта

1. Создайте Docker-сеть:
   ```bash
   docker network create my-network
   ```

2. Соберите Docker-образы для каждого сервиса:
   ```bash
   docker build -t nginx .
   docker build -t 000-user-service .
   docker build -t 001-log-service .
   docker build -t 002-robot-tester .
   ```

3. Запустите контейнеры:
   ```bash
   docker run -d --name 000-user-service --network my-network -p 5432:5432 000-user-service
   docker run -d --name 001-log-service --network my-network -p 5434:5434 001-log-service
   docker run -d --name nginx --network my-network -p 80:80 nginx
   docker run -d --name 002-robot-tester --network my-network -p 5435:5435 002-robot-tester
   ```

---

## Описание сервисов

### 000-user-service
Сервис для управления пользователями. Предоставляет API для регистрации, входа и проверки токенов.

- **Порт:** 5432
- **Dockerfile:** Использует образ Go для сборки и Alpine для запуска.
- **Основные функции:**
  - Регистрация пользователей.
  - Аутентификация пользователей.
  - Проверка токенов.

### 001-log-service
Сервис для логирования событий. Принимает логи от других сервисов и сохраняет их.

- **Порт:** 5434
- **Dockerfile:** Использует образ Go для сборки и Alpine для запуска.
- **Основные функции:**
  - Прием логов через POST-запросы.
  - Отображение логов через GET-запросы.

### nginx
Веб-сервер для маршрутизации запросов между сервисами.

- **Порт:** 80
- **Dockerfile:** Использует официальный образ NGINX.
- **Основные функции:**
  - Маршрутизация запросов.

### 002-robot-tester
Сервис для автоматического тестирования функциональности других сервисов.

- **Порт:** 5435
- **Dockerfile:** Использует образ Go для сборки и Alpine для запуска.
- **Основные функции:**
  - Автоматическая регистрация пользователей.
  - Автоматическая аутентификация.
  - Проверка токенов.

---

## API документация

### 000-user-service

#### Регистрация пользователя
- **Конечная точка:** `/register`
- **Метод:** `POST`
- **Пример запроса:**
  ```json
  {
    "username": "user1",
    "password": "1234"
  }
  ```

#### Вход пользователя
- **Конечная точка:** `/login`
- **Метод:** `POST`
- **Пример запроса:**
  ```json
  {
    "username": "user1",
    "password": "1234"
  }
  ```

#### Проверка токена
- **Конечная точка:** `/check-token`
- **Метод:** `GET`
- **Пример запроса:**
  ```bash
  curl -X GET -H "Authorization: <token>" http://localhost:5432/check-token
  ```

### 001-log-service

#### Отправка лога
- **Конечная точка:** `/log`
- **Метод:** `POST`
- **Пример запроса:**
  ```json
  {
    "message": "This is a log message"
  }
  ```

#### Получение логов
- **Конечная точка:** `/getlogs`
- **Метод:** `GET`
- **Пример запроса:**
  ```bash
  curl -X GET http://localhost:5434/getlogs
  ```

---

## Примеры запросов

### 000-user-service
```powershell
Invoke-RestMethod -Uri "http://localhost:5432/register" -Method Post -Headers @{"Content-Type"="application/json"} -Body '{"username":"user1","password":"1234"}'
Invoke-RestMethod -Uri "http://localhost:5432/login" -Method Post -Headers @{"Content-Type"="application/json"} -Body '{"username":"user1","password":"1234"}'
Invoke-RestMethod -Uri "http://localhost:5432/check-token" -Method Get -Headers @{"Authorization"="<token>"}
```

### 001-log-service
```powershell
Invoke-RestMethod -Uri "http://localhost:5434/log" -Method Post -Body (@{message = "This is a log message from PowerShell"} | ConvertTo-Json) -ContentType "application/json"
```

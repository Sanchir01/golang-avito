## 📋 Описание проекта

Сервис предназначен для управления пунктами выдачи заказов (ПВЗ), приемки товаров и ведения учета продукции. Проект разработан в рамках тестового задания и демонстрирует реализацию микросервисной архитектуры с использованием современных практик разработки на Go.

### 🌟 Основные возможности

- **Управление ПВЗ**: создание и просмотр пунктов выдачи заказов
- **Работа с приемками**: открытие и закрытие приемок товаров
- **Учет товаров**: добавление и удаление товаров в рамках приемки
- **Аутентификация и авторизация**: система ролей (модератор, сотрудник)
- **Мониторинг**: интеграция с Prometheus и Grafana

## 🏗️ Архитектура проекта

Проект построен на основе чистой архитектуры с четким разделением на слои:

```
├── cmd/               # Точки входа в приложение
├── config/            # Конфигурационные файлы
├── docs/              # Документация API (Swagger)
├── internal/          # Внутренний код приложения
│   ├── app/           # Инициализация приложения
│   ├── config/        # Структуры конфигурации
│   ├── context/       # Контекстные ключи
│   ├── feature/       # Бизнес-логика по доменам
│   │   ├── acceptance/# Функционал приемок
│   │   ├── product/   # Функционал товаров
│   │   ├── pvz/       # Функционал ПВЗ
│   │   └── user/      # Функционал пользователей
│   └── server/        # Серверная часть (HTTP, gRPC)
├── migrations/        # SQL миграции
├── pkg/               # Переиспользуемые пакеты
│   ├── db/            # Работа с базой данных
│   └── lib/           # Вспомогательные библиотеки
└── tests/             # Интеграционные тесты
```

## 🚀 Запуск проекта

### Предварительные требования

- Go 1.24+
- Docker и Docker Compose
- PostgreSQL
- Make

### Установка и запуск

1. **Клонирование репозитория**
   ```bash
   git clone https://github.com/yourusername/golang-avito.git
   cd golang-avito
   ```

2. **Запуск базы данных и сервисов мониторинга**
   ```bash
   make compose
   ```

3. **Применение миграций**
   ```bash
   make migrations-up-pg
   ```

4. **Запуск сервера**
   ```bash
   make run
   ```

5. **Запуск с автоматической перезагрузкой при изменениях**
   ```bash
   go install github.com/air-verse/air@latest
   air
   ```

## 🔌 API сервиса

### Аутентификация

#### Регистрация
```http
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "your_password",
  "role": "employee" // или "moderator"
}
```

#### Вход
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "your_password"
}
```

### Управление ПВЗ (Пункты выдачи заказов)

#### Создание ПВЗ (требуется роль модератора)
```http
POST /api/pvz
Authorization: Bearer <token>
Content-Type: application/json

{
  "registration_date": "2023-04-10T12:00:00Z",
  "city": "Москва"
}
```

#### Получение списка ПВЗ
```http
GET /api/pvz?start_date=2023-01-01T00:00:00Z&end_date=2023-12-31T23:59:59Z&page=1&limit=20
Authorization: Bearer <token>
```

### Управление приемками

#### Создание приемки
```http
POST /api/receptions
Authorization: Bearer <token>
Content-Type: application/json

{
  "pvz_id": "uuid-пвз"
}
```

#### Закрытие приемки
```http
POST /api/pvz/{pvzId}/close_last_reception
Authorization: Bearer <token>
```

### Управление товарами

#### Добавление товара
```http
POST /api/products
Authorization: Bearer <token>
Content-Type: application/json

{
  "type": "электроника",
  "acception_id": "uuid-приемки"
}
```

#### Удаление последнего товара
```http
POST /api/products/{acceptanceID}/delete_last_product
Authorization: Bearer <token>
```

## 📊 Мониторинг

Проект интегрирован с Prometheus и Grafana для мониторинга производительности:

- **Prometheus**: доступен по адресу `http://localhost:19090`
- **Grafana**: доступен по адресу `http://localhost:13000` (логин/пароль по умолчанию: admin/admin)

## 📝 Документация API

Документация API доступна через Swagger UI по адресу `http://localhost:8080/swagger/`.

## 🧪 Тестирование

### Запуск unit-тестов
```bash
go test ./...
```

### Запуск интеграционных тестов
```bash
go test -v ./tests
```

## 🔨 Полезные команды

```bash
# Сборка проекта
make build

# Запуск проекта
make run

# Создание новой миграции
make migrations-new MIGRATION_NAME=название_миграции

# Применение миграций
make migrations-up-pg

# Откат миграций
make migrations-down-pg

# Проверка статуса миграций
make migrations-status-pg

# Проверка кода линтером
make lint

# Обновление Swagger-документации
make swag
```

## 🛠️ Технологический стек

- **Язык программирования**: Go 1.24+
- **База данных**: PostgreSQL
- **API**: RESTful HTTP + gRPC
- **Документация**: Swagger
- **Контейнеризация**: Docker, Docker Compose
- **Миграции**: Goose
- **Мониторинг**: Prometheus, Grafana
- **Логирование**: slog
- **Аутентификация**: JWT

## 📈 Диаграмма базы данных

```
┌────────────┐       ┌────────────┐       ┌────────────┐
│    pvz     │       │ acceptance │       │  product   │
├────────────┤       ├────────────┤       ├────────────┤
│ id         │◄──┐   │ id         │◄──┐   │ id         │
│ city       │   │   │ pvz_id     │─┐ │   │ type       │
│ reg_date   │   └───│ status     │ │ └───│ receiving_id│
└────────────┘       │ created_at │ │     │ created_at │
                     └────────────┘ │     └────────────┘
                                    │
┌────────────┐                      │
│   users    │                      │
├────────────┤                      │
│ id         │                      │
│ email      │                      │
│ password   │                      │
│ role       │                      │
└────────────┘                      │
                                    │
```

## 👥 Разработчики

- [Sanchir01](https://github.com/Sanchir01) - Разработчик

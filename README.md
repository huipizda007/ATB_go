# 🛒 ATB_go (Fruit Market API)

RESTful API сервер, написаний на Go, для управління базою даних супермаркету. На даному етапі реалізовано мікросервіс для управління відділом свіжих фруктів (облік залишків, встановлення цін, додавання нових позицій).

Проєкт побудований з використанням чистих патернів архітектури (відділення маршрутизатора від хендлерів) та використовує сучасні інструменти екосистеми Go.

## Особливості

*   **CRUD операції:** Повне управління сутністю `Fruit` (Створення, Читання, Оновлення, Видалення).
*   **Точкові оновлення:** Окремі оптимізовані маршрути для швидкого оновлення ціни (`price`) та кількості на складі (`stock`).
*   **Чиста архітектура:** Бізнес-логіка винесена в окремий пакет `handlers`.
*   **Type-safe DB:** Використання `sqlc` для генерації безпечного Go-коду з SQL-запитів.
*   **Тестування:** Покриття HTTP-хендлерів юніт-тестами з використанням Mock-структур бази даних.

## Технологічний стек

*   **Мова:** [Go](https://go.dev/) (1.20+)
*   **Маршрутизатор:** [go-chi/chi/v5](https://github.com/go-chi/chi) (Легковагий та швидкий роутер)
*   **База даних:** PostgreSQL
*   **Драйвер БД:** [jackc/pgx/v5](https://github.com/jackc/pgx)
*   **Генерація запитів:** [sqlc](https://sqlc.dev/)
*   **Інфраструктура:** Docker

---

## Швидкий старт

### 1. Клонування репозиторію
```bash
git clone [https://github.com/huipizda007/ATB_go.git](https://github.com/huipizda007/ATB_go.git)
cd ATB_go
```

### 2. Запуск локальної бази даних (Docker)
Для роботи сервера потрібен PostgreSQL. Запустіть контейнер за допомогою цієї команди:
```bash
docker run --name atb-postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=fruits -p 5432:5432 -d postgres
```

### 3. Створення таблиці
Застосуйте початкову міграцію для створення таблиці `fruits`:
```bash
docker exec -it atb-postgres psql -U postgres -d fruits -c "CREATE TABLE fruits (id BIGSERIAL PRIMARY KEY, name TEXT NOT NULL, brand TEXT, price_per_kg NUMERIC(10, 2), stock_kg INT NOT NULL, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);"
```

### 4. Встановлення залежностей
```bash
go mod tidy
```

### 5. Запуск сервера
```bash
go run cmd/api/main.go
```
*Сервер успішно стартує на порту `:3000`.*

---

## API Документація

Базовий URL: `http://localhost:3000`

### Перевірка статусу
| Метод | Маршрут | Опис |
| :--- | :--- | :--- |
| `GET` | `/health` | Перевірка працездатності сервера |

### Управління фруктами (Fruits)
| Метод | Маршрут | Опис | Тіло запиту (JSON) |
| :--- | :--- | :--- | :--- |
| `GET` | `/fruits` | Отримати список всіх фруктів | - |
| `GET` | `/fruits/{id}` | Отримати інформацію про один фрукт | - |
| `POST` | `/fruits` | Додати новий фрукт | `{"name": "string", "brand": "string", "price_per_kg": "string", "stock_kg": int}` |
| `PATCH` | `/fruits/{id}/price` | Оновити ціну за кг | `{"price": "string"}` |
| `PATCH` | `/fruits/{id}/stock` | Оновити залишок на складі (кг) | `{"stock": int}` |
| `DELETE`| `/fruits/{id}` | Видалити фрукт з бази | - |

---

## Тестування

Проєкт містить юніт-тести для перевірки логіки HTTP-хендлерів без підключення до реальної бази даних (використовуються заглушки - mocks).

Для запуску тестів виконайте:
```bash
go test -v ./internal/server/handlers/...
```

---

## Структура проєкту

```text
ATB_go/
├── cmd/
│   └── api/
│       └── main.go           # Точка входу, конфігурація та запуск сервера
├── db/
│   ├── query/                # SQL запити для sqlc
│   ├── migration/            # Файли міграцій БД
│   └── sqlc/                 # Згенерований sqlc код для роботи з БД
├── internal/
│   └── server/
│       ├── server.go         # Налаштування маршрутизатора (Chi)
│       └── handlers/         # Бізнес-логіка та обробка HTTP-запитів
│           ├── fruit.go      # Логіка сутності Fruit
│           └── fruit_test.go # Юніт-тести
├── go.mod                    # Залежності проєкту
└── README.md                 # Документація
```
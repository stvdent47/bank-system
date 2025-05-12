# Банковская система

## Описание проекта

Банковская система представляет собой веб-приложение на языке Go, которое предоставляет функциональность для управления банковскими счетами, картами, транзакциями, кредитами и платежами. Система построена с использованием архитектуры, разделенной на слои: контроллеры, сервисы и репозитории.

## Структура базы данных

### Создание таблицы пользователей
```
create table users (
	id varchar(100) primary key,
	username varchar(50) not null unique,
	password varchar(255) not null,
	email varchar(50) not null unique
)
```

### Создание таблицы транзакций
```
create table transactions (
	id varchar(100) primary key,
	amount bigint not null,
	from_id varchar(100),
	to_id varchar(100),
	type varchar(50) not null check (type in ('deposit', 'payment', 'transfer', 'withdrawal')),
	description varchar(255) not null,
	created_at bigint not null
)
```

### Создание таблицы счетов
```
create table accounts (
	id varchar(100) primary key,
	balance bigint not null,
	user_id varchar(100) not null references users(id) on delete cascade,
	created_at bigint not null
)
```

### Создание таблицы карт
```
create table cards (
	id varchar(100) primary key,
	number bytea not null unique,
	expiration bytea not null,
	cvv varchar(255) not null,
	user_id varchar(100) not null references users(id) on delete cascade,
	account_id varchar(100) not null references accounts(id) on delete cascade,
	created_at bigint not null
)
```

### Создание таблицы платежей
```
create table payments (
	id varchar(100) primary key,
	loan_id varchar(100) not null references loans(id) on delete cascade,
	amount bigint not null,
	date bigint,
	due_date bigint not null,
	principal_part bigint not null,
	interest_part bigint not null,
	status varchar(50) not null check (status in ('new', 'paid', 'overdue')),
	is_paid boolean not null
)
```

### Создание таблицы кредитов
```
create table loans (
	id varchar(100) primary key,
	user_id varchar(100) not null references users(id) on delete cascade,
	account_id varchar(100) not null references accounts(id) on delete cascade,
	amount bigint not null,
	interest_rate double precision not null,
	term int not null,
	start_date bigint not null,
	debt bigint not null
)
```

## Архитектура приложения

### Контроллеры (Controllers) UserController

- Register - регистрация нового пользователя
- Login - аутентификация пользователя AccountController
- GetAll - получение всех счетов пользователя
- GetById - получение счета по ID
- Create - создание нового счета
- UpdateBalance - обновление баланса счета
- Delete - удаление счета
- Transfer - перевод средств между счетами CardController
- Create - создание новой карты
- GetInfo - получение информации о карте
- Pay - оплата с использованием карты LoanController
- Apply - подача заявки на кредит
- GetSchedule - получение графика платежей по кредиту AnalyticsController
- GetTransactionsAnalytics - получение аналитики по транзакциям

### Сервисы (Services) UserService

- Управление пользователями (регистрация, аутентификация) AccountService
- Создание и управление счетами
- Обновление баланса
- Перевод средств между счетами CardService
- Создание карт с шифрованием данных
- Получение информации о карте
- Обработка платежей по карте TransactionService
- Создание и управление транзакциями LoanService
- Оформление кредитов
- Формирование графика платежей
- Управление кредитной задолженностью AnalyticsService
- Анализ транзакций и финансовой активности SchedulerService
- Автоматическая проверка просроченных платежей

### Репозитории (Repositories) UserRepository

- Операции с данными пользователей в БД AccountRepository
- Операции со счетами в БД
- Обновление баланса
- Перевод средств CardRepository
- Операции с картами в БД
- Хранение зашифрованных данных карт TransactionRepository
- Операции с транзакциями в БД PaymentRepository
- Операции с платежами по кредитам
- Получение просроченных платежей LoanRepository
- Операции с кредитами в БД

## API Endpoints

### Аутентификация

- POST /register - регистрация нового пользователя
- POST /login - вход в систему

### Счета

- GET /accounts - получение всех счетов пользователя
- GET /accounts/{id} - получение счета по ID
- POST /accounts/create - создание нового счета
- PATCH /accounts/{id}/balance - обновление баланса счета
- DELETE /accounts/{id} - удаление счета
- POST /accounts/transfer - перевод средств между счетами

### Карты

- POST /cards/create - создание новой карты
- POST /cards/info - получение информации о карте
- POST /cards/pay - оплата с использованием карты

### Кредиты

- POST /loans/apply - подача заявки на кредит
- GET /loans/{id}/schedule - получение графика платежей по кредиту

### Аналитика

- POST /analytics/transactions - получение аналитики по транзакциям

## Особенности реализации

1. Безопасность данных карт :   
   - Номер карты и срок действия хранятся в зашифрованном виде (PGP)
   - CVV хранится в виде хеша (bcrypt)
2. Транзакционность :
   - Операции с деньгами выполняются в рамках транзакций БД
   - Поддержка атомарности операций
3. Аутентификация :
   - Использование JWT для аутентификации пользователей
   - Middleware для проверки токенов
4. Планировщик задач :
   - Автоматическая проверка просроченных платежей
   - Обновление статуса платежей
5. Аналитика :
   
   - Анализ транзакций пользователя
   - Формирование отчетов по финансовой активности

## Технологии
- Язык программирования : Go
- База данных : PostgreSQL
- ORM : pgx
- Маршрутизация : gorilla/mux
- Логирование : logrus
- Шифрование : PGP, bcrypt
- Аутентификация : JWT

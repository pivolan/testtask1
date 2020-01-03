# testtask1

Тестовое задание:
The main goal of this test task is a develop the application for processing the incoming requests from the 3d-party providers.  
The application must have an HTTP URL to receive incoming POST requests.  
To receive the incoming POST requests the application must have an HTTP URL endpoint.  

Technologies: Golang + Postgres.  

Requirements:  
1. Processing and saving incoming requests.  

Imagine that we have a user with the account balance.  

Example of the POST request:  
POST /your_url HTTP/1.1  
Source-Type: client  
Content-Length: 34  
Host: 127.0.0.1  
Content-Type: application/json  
{"state": "win", "amount": "10.15", "transactionId": "some generated identificator"}  

Header “Source-Type” could be in 3 types (game, server, payment). This type probably can be extended in the future.  

Possible states(win, lost):  
1. Win requests must increase the user balance  
2. Lost requests must decrease user balance.  
Each request (with the same transaction id) must be processed only once.  

The decision regarding database architecture and table structure is made to you.    

You should know that account balance can't be in a negative value.  
The application must be competitive ability.  

2. Post-processing  
Every N minutes 10 latest odd records must be canceled and balance should be corrected.  
Canceled records shouldn't be processed twice.  

## Описание  

cd ./main/  
go run main.go  

Дополнительные параметры в командной строке:  

./main -port=8098 - порт на котором слушать  
./main -dsn=host=localhost port=5432 user=postgres dbname=testtask1 sslmode=disable - строка подключения к postgres  
Для хранения баланса используем decimal тип, float для хранения баланса не используется.  
Для работы с базой данных используется gorm библиотека.  
Структура базы:  

хранение транзакций:  

    ID          string         id транзакции, уникальный ключ transation id. Предотвращение дублирования  
    OrderId     uint64         поле автоинкремент, для возможности быстрой соритровки в обратном порядке  
    CreatedAt   time.Time        
    CancelledAt *time.Time       
    Amount      decimal.Decimal может быть отрицательным значением, для упрощенного подсчета суммы в базе  
    State       StateType       Lost WIN ставка  
    UserID      uuid.UUID       Связь с пользователем, особенности библиотеки  
	User        UserBalance     Связь с пользователем  

хранение балансов пользователей:  

    ID        uuid.UUID  
    CreatedAt time.Time  
    UpdatedAt time.Time  
    DeletedAt *time.Time  
    Balance decimal.Decimal баланс пользователя  

### Процесс работы:  

Инициализируем структуру с помощью метода init()  
В нем создаем подключение к базе данных, проводим автоматическую миграцию, создаем нужные таблицы и индексы.  
Запускаем фиктуры - тестовые пользователи в базе с нулевыми балансами.  
Запускаем web server, слушаем порт на localhost  
Запускаем в фоне задачу в го рутине, которая каждые 10 минут удаляет 10 последних нечетных транзакций.  
Слушатель обрабатывает только один адрес: /my_url  
При получении запроса:  
Проверям соответствие запроса нужным параметрам, проверяем наличие нужных заголовков.  
Авторизация утрирована, берем id пользователя из заголовка и применяем к нему транзакцию.  
В транзакции проверяем результирующий баланс пользователя, если баланс меньше нуля, отменяем транзакцию.  

### Уникальность транакций:  
Блокируем таблицу на запись, создаем транзакцию, проверяем результирующий баланс, если баланс меньше нуля, отменяем транзакцию.  
Тест на условие по балансу при конкунрентной записи: TestTestTask_HandleTransactionAction_NegativeBalance  

### Крон задача:  

1.Создаем транзакцию,  
2.проверяем баланс пользователя,  
3.находим 10 транзакций,  
4.считаем баланс после удаления,  
5.удаляем  
6.проверяем баланс запросом в базу.  

### Обработка POST запроса:  

HandleTransactionAction:  
Предварительные проверки:  

1. Тип запроса: POST  
2. Заголовок Content-Type: json  
3. Source-type: Только проверяем возможные. В задаче не указано как он используется, поэтому в базе не сохранеяем, возможные варианты храним в константах.  
4. Авторизация: в задаче указан некий пользователь, про авторизацию нет данных. Используем стандартный путь в веб индустрии, пользователей много. Вместо авторизации простая заглушка, используем user_id из заголовка запроса.  
5. Проверка запроса: amount>0, state IN [lost, win]  

Создаем транзакцию  
Успешно: отдаем данные о балансе пользователя и статус - success  
Неуспешно: статус - fail, error - не пустой.  


### Тесты  

Использую минимальное количество тестов:  

postgres_actions_test.go:  
Для запуска необходимо поправить константу подключения к базе: DEFAULT_TEST_DSN  
Проверка баланса пользователя  
Отмена транзакций  
Отмена 10 последних транзакций  

web_handle_test.go:  
Для запуска необходим запущенный сервер порт 8098  
Проверка валидации заголовков  
Проверка транзакции с отрицательным баланом  
Проверка конкурентных запросов  

## Dependency  

go get "github.com/jinzhu/gorm"  
go get "github.com/lib/pq"  
go get "github.com/satori/go.uuid"  
go get "github.com/shopspring/decimal"  

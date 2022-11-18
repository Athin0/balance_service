#Сервис баланса

Использованы технологии:
- gorilla/mux  -роутер
- zap - логгер
- lib/pq -драйвер для postgres
- postgres
- go 1.18
##Структура бд
balance   -баланс пользователей  
history    -история транзакций   
reserved    -зарезервированные деньги  

##Работа приложения
Порт по умолчанию
По умолчанию использован 8080 порт и 8080 порт контейнера.  
Меняется в 
./docker-compose.yml  
- ports: *"8080:8080"*   
- SERVER_PORT: *8080*


###Запуск
```docker-compose up -d --build ```  
При запуске через ide depends работают,
при запуске через cmd иногда app запускается до postgres.
Если контейнер с app не запустился можно повторно:  
```docker-compose up -d app ```


Или чтобы ошибок точно не возникло запускать поочереди:
```docker-compose up -d postgres```  
```docker-compose up -d app```

###Остановка и удаление
```docker-compose down -v```


##Эндпоинты 
 ###Зачисление средств на счет user_id в размере value, можно добавить описание description
POST  http://localhost:8080/addIncome     {"user_id": *int*, "value": *int*, "description": *string* }

Ответ в формате:  
```{"status": "success"} ```  
При отрицательной сумме зачисления ошибки в формате:
```{"errorText": "add transaction to history query exec failed: pq: new row for relation "history" violates check constraint "history_value_check""} ```

###Резервирование средств со счета user_id в размере value
POST http://localhost:8080/addReserve     {"user_id": *int*, "order_id": *int*, "service_id": *int*,"value": *int*}  
Ответ в формате:  
```{"status": "success"} ```

Если пользователя не существует:  
```{"errorText": "user_id 0: user_id does not exist"}```

Если не хватает средств на счете:  
```{"errorText": "user_id 7: user_id has not enough balance"}```

###Услуга подтвердилась снимаем деньги у user_id в размере value  
!!если в резерве услуги не было, снимаем деньги с основного счета user_id   
POST http://localhost:8080/addExpense     {"user_id": *int*, "order_id": *int*, "service_id": *int*,"value": *int*}  
Ответ в формате:  
```{"status": "success"}```

Если пользователя не существует:  
```{"errorText": "user_id 0: user_id does not exist"}```

Если не хватает средств на счете:  
```{"errorText": "user_id 7: user_id has not enough balance"}```

###Снимаем резервацию средств со счета user_id в размере value от заказа order_id с сервиса service_id
POST http://localhost:8080/disReserve     {"user_id": *int*, "order_id": *int*, "service_id": *int*,"value": *int*}
Ответ в формате:  
```{"status": "success"}```

Если резерв не найден:  
```{"errorText": "service_id: 213, order_id:12 err: reserve does not exist"}```

###Узнает баланс пользователя user_id
GET http://localhost:8080/getBalance   {"user_id": *int*}    
Ответ в виде:  
```{"status": "success"}[{"user_id":7,"value":9024},{"user_id":89,"value":12}] ```

###Поучаем все резервации  
GET http://localhost:8080/getReserved  
Ответ в виде:
```
{"status": "success"}  
[{"id":1,"user_id":7,"service_id":212,"order_id":12,"value":12},  
{"id":2,"user_id":7,"service_id":212,"order_id":12,"value":12},  
{"id":3,"user_id":7,"service_id":212,"order_id":12,"value":12}]  
```
###Получаем балансы всех пользователй
GET http://localhost:8080/getBalances   
Ответ в виде:  
```
{"status": "success"}  
[{"user_id":7,"value":9024},  
{"user_id":89,"value":12}]
```
###Получаем историю всех транзакций
GET http://localhost:8080/getHistory  
Ответ в виде:  
```
{"status": "success"}  
[{"id":1,"user_id":7,"service_id":0,"order_id":0,"value":12,"Time":"2022-11-18T10:23:56.876456Z","description":"income salary","replenish":true},  
{"id":2,"user_id":7,"service_id":0,"order_id":0,"value":12,"Time":"2022-11-18T10:25:25.965927Z","description":"income salary","replenish":true},  
{"id":6,"user_id":7,"service_id":212,"order_id":12,"value":12,"Time":"2022-11-18T11:30:40.106175Z","description":"reserve","replenish":false}]
```
"replenish":true   --пополнение баланса  
"replenish":false  -- списание средств у пользователя

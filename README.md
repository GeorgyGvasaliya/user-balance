# user-balance
1. создать локальную базу данных в Postgresql user-balance -> таблицу users с 2мя полями - user_id, balance.
В переменную окружения DB_PASSWORD ставим пароль от бд.
Сервер запускается на localhost, 3000 порт - config.yml файл

Примеры запросов.
Имеем 4 эндпоинта
1. /get?id={id}&currency={curr}
второй параметр опционален
2. /add
в request body json
{
"user_id": {id},
"money": {money}
}
3. /withdraw аналогичен /add
4. /send
в request body json
{
"from_id": {id},
"to_id": {id},
"money":{money}
}
/get отдаёт json с user_id и balance
остальные изменяю данные, ничего в response не отдают. Только статус код.

1. Создание пользователя с балансом ![image](https://user-images.githubusercontent.com/61359396/192157509-237c66d1-3ae5-47c9-a904-e1329fe2a032.png)
2. Получение баланса пользователя ![image](https://user-images.githubusercontent.com/61359396/192157567-6b77dd91-366c-4477-b39f-433c5a57e335.png)
3. Попытка получить баланс несуществующего пользователя - NotFound + сообщение ![image](https://user-images.githubusercontent.com/61359396/192157614-e769cf47-7bf2-4d63-a016-0e2457e837c2.png)
4. Получение баланса пользователя в валюте - 24 юаня ![image](https://user-images.githubusercontent.com/61359396/192157655-9f93d672-b975-4429-bb81-88ddeb253e5d.png)
В случае невалидной валюты отдаётся BadRequest и сообщение "Wrong currency to convert". Доступные валюты брал отсюда - https://valutaomregneren.dk/data/latest.json
5. Снятие денег работает аналогично добавлению, покажу

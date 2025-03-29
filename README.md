# Matter-Poll Bot 🤖

**Matter-Poll Bot** —  это бот, который поможет Вам проводить опросы прямо в чате Mattermost.

---


## 🛠 Как запустить проект

1. **Клонировать репозиторий**
     ```bash
   git clone https://github.com/soede/matter-poll.git
   cd matter-poll
   ``` 
2. **Установить mattermost**
     ```bash
   git clone https://github.com/mattermost/docker
   cd docker
   ```
3. **Скопировать .env файл**
   ```bash
   cp ../env.example .env
   mkdir -p ./volumes/app/mattermost/{config,data,logs,plugins,client/plugins,bleve-indexes}
   ```
4. **Скопировать общий docker-compose**
   ```bash
    cp -i ../docker-compose.without-nginx.yml ./docker-compose.without-nginx.yml
   ```
5. **Запустить mattermost**
   ```bash
    docker compose -f docker-compose.yml -f docker-compose.without-nginx.yml up -d
   ```
6. **Создать токен для бота**
    1. Откройте Mattermost (http://localhost:8065/) и войдите под пользователем с правами администратора
    2. Перейдите в System Console
    3. В меню слева выберите Integration → Bot Accounts
    4. Установите Enable Bot Account Creation в true (а после нажать save)
    5. Перейдите в Main Menu → Integrations → Bot Accounts
    6. Нажмите кнопку Add Bot Account
    7. Выберите Allow this bot to post messages
    8. Нажмите Save
7. **Заполнить .env** \
   Откройте .env по адресу /mat-poll/docker/.env и заполните ``BOT_TOKEN``:
    ```plaintext
   MATTERMOST_URL=http://mattermost:8065
   BOT_TOKEN= # Вставьте сюда токен
   TARANTOOL_HOST=tarantool:3301
   TARANTOOL_PORT=3301
   TARANTOOL_USER=guest
    ``` 
8. **Запустить снова**
   ```bash
    docker compose -f docker-compose.yml -f docker-compose.without-nginx.yml down
    docker compose -f docker-compose.yml -f docker-compose.without-nginx.yml up -d --build
   ```

## Основные команды
 1️⃣ ``/guide`` – посмотреть все команды\

 2️⃣``/create Ok? | var1 | var2 | var3``, где ``Ok?`` это любой вопрос по твоему усмотрению, 
 ``var1, var2...`` – варианты ответов

 3️⃣``/vote PollID 1``, где PollID – полученный ID в ``/create`` (в след. примерах тоже)

 4️⃣``/results PollID`` – посмотреть результаты опроса

 5️⃣``/end PollID`` – завершить опрос, команда ``/results`` все еще будет актуальна, но новые голоса не принимаются
 
 6️⃣``/delete PollID`` – удалить опрос






## Сервисы
1. **Бот на Go (golang:alpine)** 
   - Работает на порту 8080.
2. **Tarantool (tarantool/tarantool:3.1.0)**
   - Написан на Go
   - Использует библиотеку ping.
   - Работает на порту 3031.
3. **Mattermost** \
   https://docs.mattermost.com/install/install-docker.html)

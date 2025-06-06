---
title: Концепт
sidebar_position: 1
---

## Проблема

Для доступа к панели администратора нужны особые права, поэтому необходимо внедрить в систему механизм аутентификации пользователей. Важно использовать бесплатный инструмент идентификации и хранить минимум персональных данных.

## Сравнение способов аутентификации

| Метод идентификации                 | Плюсы                                                                                                                                                                                                    | Минусы                                                                                                                                                       |
| ----------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| **По электронной почте**            | - Почти у всех пользователей есть email <br/> - Простая и универсальная схема (пароль + подтверждение ссылки)<br/> - Можно использовать бесплатные SMTP-сервисы (например, Gmail API, Mailgun free tier) | - Необходима реализация механизма отправки и обработки писем<br/> - Риск попадания в спам<br/> - Требуется хранить адреса                                    |
| **По номеру телефона**              | - Высокий уровень достоверности (SMS-подтверждение)<br/> - Быстрая проверка через одноразовые коды                                                                                                       | - Большинство SMS-шлюзов платные или имеют ограниченный free-tier<br/> - Нужно хранить номера телефона<br/> - Зависимость от мобильного покрытия и оператора |
| **Через Telegram (код через бота)** | - Бесплатный API Telegram Bot для отправки кода <br/> - Хранится минимум данных: только Telegram ID и временный код<br/> - Плюс верификация через учётную запись мессенджера                             | - Необходим Telegram-аккаунт<br/> - Зависимость от доступности Telegram API                                                                                  |

## Цель

Внедрить бесплатный и надежный механизм аутентификации через Telegram-бота, позволяющий администраторам входить в систему.

## Сценарии работы

### Создание пользователя в системе

1. Пользователь переходит в telegram бота и выполняет команду `/start`.
2. Бот отправляет пользователю сообщение "Привет! Для регистрации выполни команду `/register`".
3. Пользователь выполняет команду `/register`.
4. Бот отправляет в Backend API запрос на создание пользователя в системе.
5. Бот отправляет пользователю сообщение об успешной регистрации: "Ты зарегистрирован. Для повышения роли, обратись к куратору".

```mermaid
sequenceDiagram
    title Процесс регистрации

    actor User as Администратор
    participant SiteBE as Back-end
    participant Bot as Telegram-бот
    participant DB as Database

    User->>Bot: /register
    Bot->>SiteBE: POST /api/auth/save_tg_user {tg_id, username}
    SiteBE->>DB: Save user
    DB-->>SiteBE: Success
    SiteBE-->>Bot: 204 OK
    Bot-->>User: Регистрация прошла успешно
```

### Процесс авторизации пользователя в системе

1. Пользователь в веб-интерфейсе нажимает на кнопку "Войти через Telegram".
2. Бэкенд через телеграм-бота отправляет код авторизации в чат с пользователем.
3. Пользователь вводит код авторизации на сайте и попадает в панель администратора.

```mermaid
sequenceDiagram
    title Процесс авторизации

    actor User as Администратор
    participant SiteFE as Front-end
    participant SiteBE as Back-end
    participant Bot as Telegram-бот
    participant DB as Database

    User->>SiteFE: Нажатие "Войти через Telegram"
    SiteFE->>SiteBE: POST /api/auth/send_code {tg_username}
    SiteBE->>DB: Получение tg_id по tg_username
    DB-->>SiteBE: user_id или null

    alt пользователь не найден
    SiteBE-->>SiteFE: 404 "Сначала зарегистрируйтесь"
    else найден
    SiteBE->>SiteBE: Генерация 6-значного кода + expires_at
    SiteBE->>DB: Сохранение last_code, expires_at
    SiteBE->>Bot: sendMessage {chat_id: tg_id, text: code}
    SiteBE-->>SiteFE: 200 "Код отправлен, действует 10 мин"
    end

    Bot-->>User: Код подтверждения

    User->>SiteFE: Ввод кода и Submit
    SiteFE->>SiteBE: POST /api/auth/verify {tg_username, code}
    SiteBE->>DB: Получение last_code, expires_at, attempts
    DB-->>SiteBE: code_db, exp_db, attempts
    SiteBE->>SiteBE: Проверка кода
    alt валидно
    SiteBE->>SiteBE: Удалить last_code, expires_at, сброс попыток
    SiteBE->>SiteBE: сгенерировать JWT
    SiteBE-->>SiteFE: 200 {token}
    else
    SiteBE->>DB: increment attempts или invalidate code
    SiteBE-->>SiteFE: 401 "Неверный код или истёк"
    end
```

### Управление правами

#### Повышение роли пользователя

0. Проверка роли пользователя - если у пользователя есть роль куратора, то все ок. Если нет - "Ошибка: у вас недостаточно прав".
1. Куратор запрашивает список пользователей командой `/list_users`.
2. Куратор выполняет команду `/promote @user`.
3. Бот повышает роль указанного пользователя в системе до администратора.
4. Пользователю приходит уведомление в телеграм чат с ботом: "Теперь у вас есть права Администратора, ссылка на панель администратора: link".

```mermaid
sequenceDiagram
    title Процесс повышения пользователя до роли Администратор

    actor Curator as Куратор
    participant Bot as Telegram-бот
    participant DB as Database
    actor User as Пользователь

    Curator->>Bot: /list_users
    Bot->>DB: Получение пользователей
    DB-->>Bot: Список пользователей
    Bot-->>Curator: Вывод списка пользователей

    Curator->>Bot: /promote @admin_user
    Bot->>DB: Проверить роль Curator
    DB-->>Bot: role != 'curator'
    alt Роль != 'CURATOR'
        Bot-->>Curator: "Ошибка: у вас недостаточно прав"
    else Роль == 'CURATOR'
        Bot->>DB: Обновить роль @user -> 'admin'
        DB-->>Bot: Успешно
        Bot-->>Curator: "@user теперь Администратор"
        Bot-->>User: "Вам назначена роль Administrator"
    end
```

#### Понижение роли пользователя

0. Проверка роли пользователя - если у пользователя есть роль куратора, то все ок. Если нет - "Ошибка: у вас недостаточно прав".
1. Куратор запрашивает список пользователей командой `/list_admins`, выбирает понравившегося
2. Куратор выполняет команду `/demote @user`
3. Бот повышает роль указанного пользователя в системе до администратора
4. Пользователю приходит уведомление в личные сообщения "Вы потеряли статус администратора."
5. Токен доступа к панели администратора тухнет.

```mermaid
sequenceDiagram
    title Процесс повышения пользователя до роли Администратор

    actor Curator as Куратор
    participant Bot as Telegram-бот
    participant DB as Database
    actor User as Пользователь

    Curator->>Bot: /list_admins
    Bot->>DB: Получение администраторов
    DB-->>Bot: Список администраторов
    Bot-->>Curator: Вывод списка администраторов

    Curator->>Bot: /demote @user
    Bot->>DB: Проверить роль Curator
    DB-->>Bot: role != 'curator'
    alt Роль != 'CURATOR'
        Bot-->>Curator: "Ошибка: у вас недостаточно прав"
    else Роль == 'CURATOR'
        Bot->>DB: Обновить роль @user -> 'user'
        DB-->>Bot: Успешно
        Bot-->>Curator: "@user теперь обычный пользователь"
        Bot-->>Admin: "Ваша роль понижена до пользователя"
    end
```

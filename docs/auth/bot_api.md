---
title: Telegram Bot
sidebar_position: 3
---

## API

Метод для отправки сообщения пользователю в Telegram

```
POST /api/message
Authorization: Basic
{
    "telegram_user_id": "123123"
    "message": "Привет! Твой код для авторизации: ||9900||"
}
```

## Команды

- `/start`: Отправка приветственного сообщения
- `/register`: Регистрация пользователя в системе (tg_id + tg_username)
- `/list_users` (curator only): Получение списка всех пользователей
- `/list_admins` (curator only): Получение списка всех администраторов
- `/promote @user` (curator only): Повышение роли с пользователя до администратора
- `/demote @user` (curator only): Понижение роли с администратора до пользователя

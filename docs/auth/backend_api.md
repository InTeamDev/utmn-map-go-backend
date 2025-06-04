---
title: Backend
sidebar_position: 2
---

## API

Backend должен предоставить следующее API:

### Отправка кода для авторизации

```
POST /api/auth/send_code
Authorization: without
Payload:
{
    "tg_username": "@user"
}
```

**Валидация запроса:**

- Проверить, что поле `tg_username` передано и соответствует формату (начинается с `@`, длина ≤ 64).

**Логика работы:**

1. Поиск пользователя в БД по `username = tg_username`. Если не найдено – вернуть 404.
2. Проверить, когда последний раз отправлялся код: если менее 1 минуты назад – вернуть 429.
3. Сгенерировать крипто-стойкий 6-значный код с временем жизни 5 минут.
4. Сохранить код и метку `expires_at` в in-memory кэше (например, Redis) под ключом `auth_code:@user`.
5. Отправить код через Telegram Bot API.
6. Вернуть клиенту сообщение о том, что код отправлен и до какого времени он действителен.

**Нюансы:**

- Повторная отправка кода допускается не чаще одного раза в минуту.
- В случае сбоя при отправке Telegram API – логировать ошибку и выдавать 503.

**Коды ответа:**

- **200 OK**: код отправлен, в теле указан срок действия (например, `"Код действителен до 12:34:56"`).
- **400 Bad Request**: неверный или отсутствующий параметр `tg_username`.
- **404 Not Found**: пользователь с таким `tg_username` не зарегистрирован.
- **429 Too Many Requests**: запрос на отправку кода повторён ранее, чем через минуту.
- **503 Service Unavailable**: ошибка при обращении к Telegram API или кэшу.

---

### Регистрация пользователя в системе

```
POST /api/auth/save_tg_user
Authorization: Basic
Payload:
{
    "tg_id": "123123",
    "tg_username": "@user"
}
```

**Валидация запроса:**

- Проверить, что `tg_id` передан и является числом.
- Проверить, что `tg_username` передан, начинается с `@` и длина ≤ 64.

**Логика работы:**

1. Проверить в БД наличие пользователя с таким `tg_id` или `username`. Если найден – вернуть 409.
2. Вставить новую запись в `tg_users` с полями `tg_id`, `username`, `created_at`, `updated_at`.

**Коды ответа:**

- **204 No Content**: пользователь успешно создан.
- **400 Bad Request**: отсутствуют или некорректны `tg_id`/`tg_username`.
- **401 Unauthorized**: неверная Basic-авторизация сервис-to-service.
- **409 Conflict**: такой `tg_id` или `tg_username` уже зарегистрирован.

---

### Проверка кода авторизации

```
POST /api/auth/verify
Authorization: without
Payload:
{
    "tg_username": "@user",
    "code": "423142"
}
```

**Валидация запроса:**

- Проверить, что `tg_username` передан и в формате `@...`.
- Проверить, что `code` — строка из 6 цифр.

**Логика работы:**

1. Найти `user_id` в `tg_users` по `username = tg_username`. Если не найден → 404.
2. Получить из кэша запись `auth_code:@user` (код, `expires_at`, `attempts`). Если нет записи → вернуть 400.
3. Если `now() > expires_at` → удалить запись из кэша и вернуть 400.
4. Если `attempts ≥ 3` → вернуть 429.
5. Если `code` не совпал с сохранённым → увеличить `attempts` в кэше и вернуть 401.
6. При совпадении:
   - Удалить запись из кэша.
   - Сгенерировать JWT с `exp = now() + 30 min` и включить `user_id` + роли.
   - Вернуть токен в ответе.

**Нюансы:**

- Лимит попыток проверки — 3 раза, затем блокировка до истечения TTL.
- После успешной верификации код удаляется.

**Коды ответа:**

- **200 OK**: код верный, в теле `{ "token": "..." }`.
- **400 Bad Request**: код не запрашивался или просрочен, либо неверный формат полей.
- **401 Unauthorized**: код неверный (увеличивается счётчик попыток).
- **404 Not Found**: пользователь с указанным `tg_username` не найден.
- **429 Too Many Requests**: превышен лимит неверных попыток ввода кода (≥ 3).

---

### Обновление токена

```
POST /api/auth/refresh
Authorization: without (Refresh-Token в HttpOnly cookie)
Payload:
{}
```

**Валидация запроса:**

- Проверить наличие HttpOnly cookie `refresh_token` (или поля в теле) и что оно не пусто.

**Логика работы:**

1. Извлечь `refresh_token` из cookie. Если нет — вернуть 401.
2. Декодировать и верифицировать подпись JWT (HS256, секрет `JWT_SECRET`). При ошибке — 401.
3. Проверить в payload поля `exp` (не просрочен) и `jti` (UUID).
4. В таблице `refresh_tokens` найти запись по `jti_rt = jti`:
   - Если нет записи или `revoked = true` — вернуть 401.
   - Если `now() > expires_at` — вернуть 401.
5. Сгенерировать новую пару токенов:
   - **Access Token** (exp через 15 мин, новый `jti`).
   - **Refresh Token** (exp через 30 дн, новый `jti_rt`).
6. В таблице `refresh_tokens` пометить старую запись `revoked = true` и вставить новую с новым `jti_rt`.
7. Установить новый `refresh_token` в HttpOnly cookie и вернуть в теле `{ "access_token": "…" }`.

**Нюансы:**

- Ротация Refresh Token: старый помечается `revoked`, новый выдаётся при каждом обновлении.
- При частых обновлениях можно ограничить частоту (например, не чаще чем раз в 5 минут).

**Коды ответа:**

- **200 OK**: выдан новый access_token и установлен новый refresh_token.
- **400 Bad Request**: отсутствует или некорректен refresh_token.
- **401 Unauthorized**: подпись не прошла, токен просрочен или отозван.
- **429 Too Many Requests**: превышена частота обновлений (опционально).
- **500 Internal Server Error**: внутренняя ошибка сервиса.

---

### Logout

```
POST /api/auth/logout
Authorization: Bearer <access_token> (Refresh-Token в HttpOnly cookie)
Payload:
{}
```

**Валидация запроса:**

- Проверить наличие заголовка `Authorization: Bearer` с access_token.

**Логика работы:**

1. Верифицировать `access_token`. При ошибке — вернуть 401.
2. Из payload извлечь `jti` access-токена и `sub` (user_id).
3. Вставить `jti` access-токена в таблицу `jwt_blacklist`.
4. Извлечь из cookie текущий `refresh_token`, декодировать, извлечь `jti_rt`.
5. В таблице `refresh_tokens` найти запись по `jti_rt`:
   - Если найдена — пометить `revoked = true`.
6. Очистить cookie `refresh_token` (установить истёкшую дату).
7. Вернуть **204 No Content**.

**Нюансы:**

- Logout одновременно отзывает и access-, и refresh-токен.
- Даже если клиент не отправил refresh_token (cookie удалена раньше) — отзыв access-токена всё равно проводится.

**Коды ответа:**

- **204 No Content**: выход выполнен — токены отозваны.
- **401 Unauthorized**: неверный или просроченный access_token.
- **500 Internal Server Error**: внутренняя ошибка при работе с базой.

## Middleware для Basic авторизации

Для s2s (service-to-service) взаимодействия, используется Basic авторизация, где login - client_id, а пароль - access_token.

**Логика работы:**

1. Из Authorization заголовка извлекается токен. Если токена нет – возвращать **401 Unauthorized.**
2. Декодировать токен:
   - Убрать префикс `Basic`, декодировать оставшуюся часть из Base64 в строку `clientid:accesstoken`.
   - Если формат не соответствует - **401 Unauthorized.**
3. Сравнить с auth-config:
   - Сравнить cliendid и accesstoken с теми, что лежат в auth-config.
   - Если формат не соответствует - **401 Unauthorized.**

## Middleware для JWT авторизации

Для взаимодействия User-Backend, используется авторизация по JWT токену.

**Логика работы:**

1. Из Authorization заголовка извлекается токен. Если токена нет – возвращать **401 Unauthorized.**
2. Валидация подписи и декодирование (HS256):
   - Убрать префикс `Bearer`, получив чистый JWT.
   - Разбить на три части: `header.payload.signature`.
   - Посчитать HMAC-SHA256 от `header.payload` с секретом `JWT_SECRET`.
   - Base64URL-кодировать результат и сравнить с переданной `signature`.
   - Если подпись не совпадает – возвращать **401 Unauthorized.**
3. Проверка claims полей:
   - `exp (expiration)` должен быть больше текущего времени; иначе – **401 Unauthorized.**
   - `nbf (not before)` должен быть меньше или равен текущего времени; иначе – **401 Unauthorized.**
4. Проверка отзыва токена:
   - Из payload извлечь `jti`.
   - Проверить в таблице `jwt_blacklist`: если есть запись с таким `jti` – возвращать **401 Unauthorized.**

## Выпуск JWT-токена

Для авторизации пользователя в системе, нужно предоставить JWT-токен авторизации.

Для обновления токена авторизации без повторного ввода кода из telegram, нужно обновлять токен.
Для этого выпускается пара `access + refresh token`.

**Логика работы:**

JWT токен состоит из 3 частей: `<header>.<payload>.<signature>`

Генерация Access Token:

1. Формирование header:
   - Определяются `{"alg": "HS256", "typ": "JWT"}`
   - Json кодируется в Base64
2. Формирование payload:
   - Определяется:
   ```json
   {
     "sub": "user_id",
     "roles": ["user", "admin"],
     "iat": now(),
     "exp": now() + 30days,
     "jti": uuid() // uuid чтобы блокнуть токен
   }
   ```
   - Json кодируется в Base64
3. Формирование подписи:
   - Формирование строки `signingInput = "<header>.<payload>"`
   - Вычисляем HMAC-SHA256 `signature = HMAC-SHA256(signingInput, JWT_SECRET)`
   - Результат кодируется в Base64
4. Конкатенация частей:
   - `accessToken = <header>.<payload>.<signature>`

Генерация Refresh Token:

1. Формирование header:
   - Определяются `{"alg": "HS256", "typ": "JWT"}`
   - Json кодируется в Base64
2. Формирование payload:
   - Определяется:
   ```json
   {
     "sub": "user_id",
     "iat": now(),
     "exp": now() + 15min,
     "jti": uuid() // uuid чтобы блокнуть токен
   }
   ```
   - Json кодируется в Base64
3. Формирование подписи:
   - Формирование строки `signingInputRT = "<header_rt>.<payload_rt>"`
   - Вычисляем HMAC-SHA256 `signatureRT = HMAC-SHA256(signingInputRT, JWT_SECRET)`
   - Результат кодируется в Base64
4. Конкатенация частей:
   - `refreshToken = <header_rt>.<payload_rt>.<signature_rt>`
5. Сохранение Refresh Token (`jti_rt`, `iat`, `exp`, `user_id`) в таблицу `refresh_tokens`.
6. Возвращаем клиенту пару токенов:
   ```json
   {
     "access_token": "<accessToken>",
     "refresh_token": "<refreshToken>"
   }
   ```

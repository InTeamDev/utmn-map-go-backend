-- name: GetAdminByUsername :one
SELECT 
    id,
    username,
    password_hash,
    email,
    full_name,
    is_active,
    last_login
FROM administrators
WHERE username = @username;

-- name: CreateAdmin :one
INSERT INTO administrators (
    username,
    password_hash,
    email,
    full_name
) VALUES (
    @username,
    @password_hash,
    @email,
    @full_name
)
RETURNING id, username, email, full_name, created_at;

-- name: UpdateAdminLastLogin :exec
UPDATE administrators
SET last_login = CURRENT_TIMESTAMP
WHERE id = @id;

-- name: UpdateAdminPassword :exec
UPDATE administrators
SET 
    password_hash = @new_password_hash,
    password_reset_token = NULL,
    token_expires_at = NULL
WHERE id = @id;

-- name: SetPasswordResetToken :exec
UPDATE administrators
SET 
    password_reset_token = @token,
    token_expires_at = @expires_at
WHERE email = @email;

-- name: GetAdminByResetToken :one
SELECT 
    id,
    username,
    email,
    token_expires_at
FROM administrators
WHERE password_reset_token = @token;

-- name: GetAllAdmins :many
SELECT 
    id,
    username,
    email,
    full_name,
    created_at,
    last_login,
    is_active
FROM administrators
ORDER BY created_at DESC;

-- name: UpdateAdminStatus :exec
UPDATE administrators
SET is_active = @is_active
WHERE id = @id;

-- name: DeleteAdmin :exec
DELETE FROM administrators
WHERE id = @id;

-- name: UpdateAdminProfile :exec
UPDATE administrators
SET 
    email = @email,
    full_name = @full_name
WHERE id = @id;

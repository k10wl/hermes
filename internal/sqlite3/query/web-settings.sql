-- name: GetWebSettings :one
SELECT * FROM web_settings LIMIT 1;

-- name: UpdateWebSettings :exec
UPDATE web_settings SET dark_mode = ?, initted = ?;

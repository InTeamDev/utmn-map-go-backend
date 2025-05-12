-- name: CreateConnection :one
INSERT INTO connections (from_id, to_id, weight)
VALUES (@from_id, @to_id, @weight)
RETURNING *;

-- name: CreateIntersection :one
INSERT INTO intersections (id, x, y, floor_id)
VALUES (@id, @x, @y, @floor_id)
RETURNING *;

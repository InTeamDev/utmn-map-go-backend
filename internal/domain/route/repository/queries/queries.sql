-- name: CreateConnection :one
INSERT INTO connections (from_id, to_id, weight)
VALUES (@from_id, @to_id, @weight)
RETURNING *;

-- name: CreateIntersection :one
INSERT INTO intersections (id, x, y, floor_id)
VALUES (@id, @x, @y, @floor_id)
RETURNING *;

-- name: GetConnections :many
SELECT c.from_id, c.to_id, c.weight 
FROM connections c
WHERE EXISTS (
    SELECT 1 FROM intersections i
    JOIN floors f ON i.floor_id = f.id
    WHERE (i.id = c.from_id OR i.id = c.to_id)
    AND f.building_id = $1
);
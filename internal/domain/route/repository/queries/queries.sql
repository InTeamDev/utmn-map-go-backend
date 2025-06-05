-- name: CreateConnection :one
INSERT INTO connections (from_id, to_id, weight)
VALUES (@from_id, @to_id, @weight)
ON CONFLICT (from_id, to_id) DO UPDATE
  SET weight = EXCLUDED.weight
RETURNING *;

-- name: CreateIntersection :one
INSERT INTO intersections (id, x, y, floor_id)
VALUES (@id, @x, @y, @floor_id)
ON CONFLICT (id) DO UPDATE
  SET x = EXCLUDED.x,
      y = EXCLUDED.y,
      floor_id = EXCLUDED.floor_id
RETURNING *;

-- name: DeleteIntersectionConnections :exec
DELETE FROM connections
WHERE from_id = @intersection_id OR to_id = @intersection_id;

-- name: DeleteIntersection :exec
DELETE FROM intersections i
USING floors f, buildings b
WHERE i.id = @intersection_id
  AND i.floor_id = f.id
  AND f.building_id = b.id
  AND b.id = @building_id;

-- name: GetConnections :many
SELECT c.from_id, c.to_id, c.weight 
FROM connections c
WHERE EXISTS (
    SELECT 1 FROM intersections i
    JOIN floors f ON i.floor_id = f.id
    WHERE (i.id = c.from_id OR i.id = c.to_id)
    AND f.building_id = $1
);

-- name: GetIntersections :many
SELECT 
    i.id,
    i.x,
    i.y,
    f.id AS floor_id,
    b.id AS building_id
FROM intersections i
JOIN floors f ON i.floor_id = f.id
JOIN buildings b ON f.building_id = b.id
WHERE b.id = @building_id::uuid;

-- name: GetNodeBuilding :one
SELECT f.building_id
FROM intersections i
JOIN floors f ON i.floor_id = f.id
WHERE i.id = @node_id::uuid
UNION
SELECT f.building_id
FROM doors d
JOIN objects o ON d.object_id = o.id
JOIN floors f ON o.floor_id = f.id
WHERE d.id = @node_id::uuid
LIMIT 1;

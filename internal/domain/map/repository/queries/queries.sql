-- name: GetBuildings :many
SELECT 
    b.id, 
    b.name,
    b.address
FROM buildings b;

-- name: GetBuildingByID :one
SELECT 
    b.id, 
    b.name,
    b.address
FROM buildings b
WHERE b.id = @id::uuid;

-- name: ListPolygonPointsByPolygonID :many
SELECT
  polygon_id,
  point_order AS order,
  x,
  y
FROM floor_polygon_points
WHERE polygon_id = $1
ORDER BY point_order;

-- name: GetFloorsByBuilding :many
SELECT 
    f.id, 
    f.name,
    f.alias,
    f.building_id
FROM floors f 
WHERE f.building_id = @building_id::uuid ORDER BY name asc;

-- name: GetDoorsByBuilding :many
SELECT
    d.id,
    d.x,
    d.y,
    d.width,
    d.height,
    b.id AS building_id,
    o.id AS object_id
FROM doors d
JOIN objects o ON d.object_id = o.id
JOIN floors f ON f.id = o.floor_id
JOIN buildings b ON f.building_id = b.id
WHERE b.id = @building_id::uuid;

-- name: GetObjectTypes :many
SELECT DISTINCT ot.*
FROM object_types ot
ORDER BY ot.id;

-- name: GetObjectsByBuilding :many
SELECT 
    o.id, 
    o.name, 
    o.alias, 
    o.description, 
    o.x, 
    o.y, 
    o.width, 
    o.height, 
    ot.id AS object_type, 
    f.id AS floor_id, 
    f.name AS floor_name, 
    b.id AS building_id, 
    b.name AS building_name
FROM objects o
JOIN floors f ON o.floor_id = f.id
JOIN buildings b ON f.building_id = b.id
JOIN object_types ot ON o.object_type_id = ot.id
WHERE b.id = @building_id::uuid;

-- name: GetDoorsByObjectIDs :many
SELECT 
    d.id, 
    d.x, 
    d.y, 
    d.width, 
    d.height,
    d.object_id
FROM doors d
WHERE d.object_id = ANY(@object_ids::uuid[]);

-- name: GetObjectByID :one
SELECT 
    o.*, 
    f.*, 
    (
        SELECT json_agg(json_build_object(
            'id', d.id,
            'x', d.x,
            'y', d.y,
            'width', d.width,
            'height', d.height
        ))
        FROM doors d
        WHERE d.object_id = o.id
    ) AS doors
FROM objects o
JOIN floors f ON o.floor_id = f.id
WHERE o.id = $1;

-- name: UpdateObject :one
UPDATE objects
SET name = COALESCE(sqlc.narg('name'), name),
    alias = COALESCE(sqlc.narg('alias'), alias),
    description = COALESCE(sqlc.narg('description'), description),
    x = COALESCE(sqlc.narg('x'), x),
    y = COALESCE(sqlc.narg('y'), y),
    width = COALESCE(sqlc.narg('width'), width),
    height = COALESCE(sqlc.narg('height'), height),
    object_type_id = COALESCE(sqlc.narg('object_type_id'), object_type_id)
WHERE id = @id
RETURNING *;

-- name: CreateObject :one
INSERT INTO objects (
    id,
    floor_id,
    name,
    alias,
    description,
    x,
    y,
    width,
    height,
    object_type_id
) VALUES (
    @id,
    @floor_id,
    @name,
    @alias,
    @description,
    @x,
    @y,
    @width,
    @height,
    @object_type_id
)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    alias = EXCLUDED.alias,
    description = EXCLUDED.description,
    x = EXCLUDED.x,
    y = EXCLUDED.y,
    width = EXCLUDED.width,
    height = EXCLUDED.height,
    object_type_id = EXCLUDED.object_type_id,
    floor_id = EXCLUDED.floor_id
RETURNING *;

-- name: DeleteObject :exec
DELETE FROM objects
WHERE id = @id::uuid;

-- name: GetPolygonByID :one
SELECT id, floor_id, label, z_index
FROM floor_polygons
WHERE id = $1;

-- name: GetFloorByID :one
SELECT * FROM floors
WHERE id = @id::uuid;

-- name: GetObjectTypeByID :one
SELECT * FROM object_types
WHERE id = @id::int;

-- name: GetObjectTypeByName :one
SELECT * FROM object_types
WHERE name = @name::VARCHAR(50);

-- name: GetFloorBackground :many
SELECT 
    fp.id, 
    fp.label, 
    fp.z_index, 
    json_agg(
        json_build_object(
            'order', fpp.point_order,
            'x', fpp.x,
            'y', fpp.y
        ) ORDER BY fpp.point_order
    ) AS points
FROM floor_polygons fp
JOIN floor_polygon_points fpp ON fp.id = fpp.polygon_id
WHERE fp.floor_id = @floor_id::uuid
GROUP BY fp.id, fp.label, fp.z_index
ORDER BY fp.z_index;

-- name: CreateBuilding :one
INSERT INTO buildings (id, name, address)
VALUES ($1, $2, $3)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    address = EXCLUDED.address
RETURNING *;

-- name: DeleteBuilding :exec
DELETE FROM buildings
WHERE id = $1;

-- name: UpdateBuilding :one
UPDATE buildings
SET name = COALESCE(sqlc.narg('name'), name),
address = COALESCE(sqlc.narg('address'), address)
WHERE id = sqlc.arg('id')::uuid
RETURNING id, name, address;

-- name: CreateFloor :exec
INSERT INTO floors (id, name, alias, building_id)
VALUES (@id::uuid, @name, @alias, @building_id::uuid)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    alias = EXCLUDED.alias,
    building_id = EXCLUDED.building_id;

-- name: CreateDoor :one
INSERT INTO doors (id, x, y, width, height, object_id)
VALUES (@id::uuid, @x, @y, @width, @height, @object_id::uuid)
ON CONFLICT (id) DO UPDATE SET
    x = EXCLUDED.x,
    y = EXCLUDED.y,
    width = EXCLUDED.width,
    height = EXCLUDED.height,
    object_id = EXCLUDED.object_id
RETURNING id, x, y, width, height, object_id;

-- name: GetDoor :one
SELECT 
    d.id,
    d.x,
    d.y,
    d.width,
    d.height,
    d.object_id
FROM 
    doors d
WHERE 
    d.id = @doorID::uuid
    AND EXISTS (
        SELECT 1 FROM objects o
        JOIN floors f ON o.floor_id = f.id
        WHERE o.id = d.object_id
          AND f.building_id = @buildingID::uuid
    );

-- name: UpdateDoor :one
UPDATE doors
SET 
    x = COALESCE($1, x),
    y = COALESCE($2, y),
    width = COALESCE($3, width),
    height = COALESCE($4, height),
    object_id = COALESCE($5, object_id)
WHERE 
    id = @door_id::uuid
    AND EXISTS (
        SELECT 1 FROM objects o
        JOIN floors f ON o.floor_id = f.id
        WHERE o.id = doors.object_id
          AND f.building_id = @building_id::uuid
    )
RETURNING 
    id,
    x,
    y,
    width,
    height,
    object_id;

-- name: CreatePolygon :one
INSERT INTO floor_polygons (id, floor_id, label, z_index)
VALUES (@id::uuid, @floor_id::uuid, @label, @z_index)
ON CONFLICT (id) DO UPDATE SET
    floor_id = EXCLUDED.floor_id,
    label = EXCLUDED.label,
    z_index = EXCLUDED.z_index
RETURNING id, floor_id, label, z_index;

-- name: CreatePolygonPoint :one
INSERT INTO floor_polygon_points (polygon_id, point_order, x, y)
VALUES (@polygon_id::uuid, @point_order, @x, @y)
ON CONFLICT (polygon_id, point_order) DO UPDATE SET
    x = EXCLUDED.x,
    y = EXCLUDED.y
RETURNING polygon_id, point_order, x, y;

-- name: DeletePolygonPoints :exec
DELETE FROM floor_polygon_points
WHERE polygon_id = $1 AND point_order = ANY($2::int[]);

-- name: GetDoorFloorPairs :many
SELECT
  d.id       AS door_id,
  o.floor_id AS floor_id
FROM doors d
JOIN objects o ON d.object_id = o.id;

-- name: GetObjectDoorPairs :many
SELECT
  d.id AS door_id,
  o.id AS object_id
FROM doors d
JOIN objects o ON d.object_id = o.id;

-- name: GetPolygonsByFloorID :many
SELECT id, floor_id, label, z_index
FROM floor_polygons
WHERE floor_id = $1;

-- name: ChangePolygon :exec
UPDATE floor_polygons
SET
  label = COALESCE(sqlc.narg('label')::text, label),
  z_index = COALESCE(sqlc.narg('z_index')::int, z_index)
WHERE id = @id::uuid;
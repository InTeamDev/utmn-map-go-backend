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

-- name: GetFloorsByBuilding :many
SELECT 
    f.id, 
    f.name,
    f.alias,
    f.building_id
FROM floors f 
WHERE f.building_id = @building_id::uuid;

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
    od.object_id
FROM doors d
JOIN object_doors od ON d.id = od.door_id
WHERE od.object_id = ANY(@object_ids::uuid[]);

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

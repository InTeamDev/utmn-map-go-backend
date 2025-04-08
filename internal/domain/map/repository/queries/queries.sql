-- name: GetBuildings :many
SELECT 
    b.id, 
    b.name,
    b.address
FROM buildings b;

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
JOIN objects o ON o.object_type_id = ot.id
JOIN floors f ON o.floor_id = f.id
WHERE f.building_id = @building_id::uuid;

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
    ot.name AS object_type, 
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

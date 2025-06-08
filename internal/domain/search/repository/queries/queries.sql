-- name: Search :many
SELECT DISTINCT
    o.id AS object_id,
    o.object_type_id AS object_type_id,
    d.id AS door_id,
    CONCAT(o.name, ' (', f.name, ')') AS preview
FROM
    objects o
    JOIN floors f ON f.id = o.floor_id
    JOIN buildings b ON b.id = f.building_id
    JOIN doors d ON d.object_id = o.id
    JOIN connections c ON (
        c.from_id = d.id
        OR c.to_id = d.id
    )
WHERE
    b.id = $1
    AND (
        $2 = ''
        OR o.name ILIKE '%' || $2 || '%'
        OR o.alias ILIKE '%' || $2 || '%'
        OR o.description ILIKE '%' || $2 || '%'
    )
ORDER BY preview, o.id;
CREATE TABLE floor_polygons (
    id UUID PRIMARY KEY,
    floor_id UUID NOT NULL,
    label VARCHAR(255), -- например: "стена", "пол", "зона разделения" и т.п.
    z_index INT DEFAULT 0, -- порядок отрисовки (элементы с меньшим z_index отрисовываются раньше)
    FOREIGN KEY (floor_id) REFERENCES floors (id) ON DELETE CASCADE
);

CREATE TABLE floor_polygon_points (
    polygon_id UUID NOT NULL,
    point_order INT NOT NULL,
    x FLOAT NOT NULL,
    y FLOAT NOT NULL,
    FOREIGN KEY (polygon_id) REFERENCES floor_polygons (id) ON DELETE CASCADE,
    PRIMARY KEY (polygon_id, point_order)
);

---- create above / drop below ----

DROP TABLE IF EXISTS floor_polygon_points;
DROP TABLE IF EXISTS floor_polygons;

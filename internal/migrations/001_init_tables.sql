CREATE TABLE buildings (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL
);

CREATE TABLE floors (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    alias VARCHAR(255) NOT NULL,
    building_id UUID NOT NULL,
    FOREIGN KEY (building_id) REFERENCES buildings (id) ON DELETE CASCADE
);

CREATE TABLE object_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    alias VARCHAR(50) UNIQUE NOT NULL
);

INSERT INTO object_types (name, alias) VALUES
    ('cabinet', 'Аудитория'),
    ('department', 'Кафедра'),
    ('man-toilet', 'Мужской туалет'),
    ('woman-toilet', 'Женский туалет'),
    ('stair', 'Лестница'),
    ('wardrobe', 'Гардероб'),
    ('gym', 'Спортзал'),
    ('cafe', 'Кафетерий'),
    ('canteen', 'Столовая'),
    ('chill-zone', 'Зона отдыха');

CREATE TABLE objects (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    alias VARCHAR(255) NOT NULL,
    description TEXT,
    x FLOAT NOT NULL,
    y FLOAT NOT NULL,
    width FLOAT NOT NULL,
    height FLOAT NOT NULL,
    object_type_id INT NOT NULL,
    floor_id UUID NOT NULL,
    FOREIGN KEY (object_type_id) REFERENCES object_types (id) ON DELETE RESTRICT,
    FOREIGN KEY (floor_id) REFERENCES floors (id) ON DELETE CASCADE
);

CREATE TABLE doors (
    id UUID PRIMARY KEY,
    x FLOAT NOT NULL,
    y FLOAT NOT NULL,
    width FLOAT NOT NULL,
    height FLOAT NOT NULL,
    object_id UUID NOT NULL,
    FOREIGN KEY (object_id) REFERENCES objects (id) ON DELETE CASCADE
);

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

CREATE TABLE intersections (
    id UUID PRIMARY KEY,
    x DOUBLE PRECISION NOT NULL,
    y DOUBLE PRECISION NOT NULL,
    floor_id UUID NOT NULL,
    FOREIGN KEY (floor_id) REFERENCES floors (id) ON DELETE CASCADE
);

CREATE TABLE connections (
    from_id UUID NOT NULL,
    to_id UUID NOT NULL,
    weight DOUBLE PRECISION NOT NULL CHECK (weight >= 0),
    PRIMARY KEY (from_id, to_id)
);

---- create above / drop below ----

DROP TABLE doors;

DROP TABLE objects;

DELETE FROM object_types;

DROP TABLE object_types;

DROP TABLE floors;

DROP TABLE buildings;

DROP TABLE floor_polygon_points;

DROP TABLE floor_polygons;

DROP TABLE intersections;

DROP TABLE connections;
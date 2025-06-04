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
    FOREIGN KEY (building_id) REFERENCES buildings(id) ON DELETE CASCADE
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
    ('cafe', 'Кафе'),
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
    FOREIGN KEY (object_type_id) REFERENCES object_types(id) ON DELETE RESTRICT,
    FOREIGN KEY (floor_id) REFERENCES floors(id) ON DELETE CASCADE
);

CREATE TABLE doors (
    id UUID PRIMARY KEY,
    x FLOAT NOT NULL,
    y FLOAT NOT NULL,
    width FLOAT NOT NULL,
    height FLOAT NOT NULL,
    object_id UUID NOT NULL,
    FOREIGN KEY (object_id) REFERENCES objects(id) ON DELETE CASCADE
);

---- create above / drop below ----

DROP TABLE doors;
DROP TABLE objects;
DELETE FROM object_types;
DROP TABLE object_types;
DROP TABLE floors;
DROP TABLE buildings;

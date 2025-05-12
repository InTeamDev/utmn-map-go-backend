CREATE TABLE intersections (
  id UUID PRIMARY KEY,
  x  DOUBLE PRECISION NOT NULL,
  y  DOUBLE PRECISION NOT NULL,
  floor_id UUID NOT NULL,
  FOREIGN KEY (floor_id) REFERENCES floors(id) ON DELETE CASCADE
);

CREATE TABLE connections (
  from_id UUID NOT NULL,
  to_id   UUID NOT NULL,
  weight  DOUBLE PRECISION NOT NULL CHECK (weight >= 0),
  PRIMARY KEY (from_id, to_id)
);

---- create above / drop below ----

DROP TABLE intersections;
DROP TABLE connections;
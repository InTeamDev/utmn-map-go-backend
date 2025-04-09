-- 1 этаж — Лестница 1 (южная стена)
INSERT INTO objects (id, name, alias, x, y, width, height, object_type_id, floor_id)
VALUES ('a98ef279-e4a7-44e4-b093-71221e3b8e4f', 'Лестница 1', 'Лестница у южной стены', 247.1, 1074.5, 66.8, 92.1,
  (SELECT id FROM object_types WHERE name = 'stair'), 'd33d56e3-aca8-4b62-b34a-9e3882276f75');

-- 1 этаж — Лестница 2 (центральная левая)
INSERT INTO objects (id, name, alias, x, y, width, height, object_type_id, floor_id)
VALUES ('b71ae235-8b14-4175-9607-3068d49a3c8a', 'Лестница 2', 'Центральная левая', 753.3, 286.7, 92.1, 67.3,
  (SELECT id FROM object_types WHERE name = 'stair'), 'd33d56e3-aca8-4b62-b34a-9e3882276f75');

-- 1 этаж — Лестница 3 (центральная правая)
INSERT INTO objects (id, name, alias, x, y, width, height, object_type_id, floor_id)
VALUES ('87d3a8f5-2795-4530-998e-20bbbc48240e', 'Лестница 3', 'Центральная правая', 964.8, 286.7, 92.1, 67.3,
  (SELECT id FROM object_types WHERE name = 'stair'), 'd33d56e3-aca8-4b62-b34a-9e3882276f75');

-- 1 этаж — Лестница 4 (северо-восточная секция)
INSERT INTO objects (id, name, alias, x, y, width, height, object_type_id, floor_id)
VALUES ('20dcdfbd-2cd7-4861-bb0c-3be1ab1504a1', 'Лестница 4', 'Северо-восточная секция', 865.1, 1106.8, 41.6, 90.9,
  (SELECT id FROM object_types WHERE name = 'stair'), 'd33d56e3-aca8-4b62-b34a-9e3882276f75');

-- 1 этаж — Лестница 5 (правый верхний угол)
INSERT INTO objects (id, name, alias, x, y, width, height, object_type_id, floor_id)
VALUES ('2f04a0ae-69b3-4534-9061-97c340e261f3', 'Лестница 5', 'Правый верхний угол', 1816.4, 903, 87.6, 107.5,
  (SELECT id FROM object_types WHERE name = 'stair'), 'd33d56e3-aca8-4b62-b34a-9e3882276f75');

-- 1 этаж — Лестница 6 (входная)
INSERT INTO objects (id, name, alias, x, y, width, height, object_type_id, floor_id)
VALUES ('82a1a2b2-356e-4f61-a07d-3e6f48a6b2a2', 'Лестница 6', 'Входная лестница', 1589.2, 57.1, 67.3, 92.1,
  (SELECT id FROM object_types WHERE name = 'stair'), 'd33d56e3-aca8-4b62-b34a-9e3882276f75');

-- 1 этаж — Лестница 7 (центр юг)
INSERT INTO objects (id, name, alias, x, y, width, height, object_type_id, floor_id)
VALUES ('9b728658-f5fd-4f31-b872-90f60e1eac36', 'Лестница 7', 'Центр юг', 442.4, 915, 41.6, 73.6,
  (SELECT id FROM object_types WHERE name = 'stair'), 'd33d56e3-aca8-4b62-b34a-9e3882276f75');

-- 1 этаж — Лестница 8 (левый юг)
INSERT INTO objects (id, name, alias, x, y, width, height, object_type_id, floor_id)
VALUES ('fd8d2a93-55e3-4a5f-87a0-c9a12f0021cb', 'Лестница 8', 'Левый юг', 770.8, 838.4, 73.6, 41.6,
  (SELECT id FROM object_types WHERE name = 'stair'), 'd33d56e3-aca8-4b62-b34a-9e3882276f75');

-- 1 этаж — Лестница 9 (правый вход)
INSERT INTO objects (id, name, alias, x, y, width, height, object_type_id, floor_id)
VALUES ('beed2307-33c6-4d10-a416-73d31edb5aa9', 'Лестница 9', 'Правый вход', 733, 370.3, 30, 73.6,
  (SELECT id FROM object_types WHERE name = 'stair'), 'd33d56e3-aca8-4b62-b34a-9e3882276f75');

-- 1 этаж — Лестница 10 (западный выход)
INSERT INTO objects (id, name, alias, x, y, width, height, object_type_id, floor_id)
VALUES ('8b4c9b42-80fd-4a85-9d52-ef78a948f2e4', 'Лестница 10', 'Западный выход', 1361.9, 712.7, 41.6, 73.6,
  (SELECT id FROM object_types WHERE name = 'stair'), 'd33d56e3-aca8-4b62-b34a-9e3882276f75');

-- 2 этаж — Лестница 1 (южная стена)
INSERT INTO objects (id, name, alias, x, y, width, height, object_type_id, floor_id)
VALUES ('2ae8c611-2321-46d6-933f-bb5e17c6bb38', 'Лестница 1', 'Лестница у южной стены', 247.1, 1074.5, 66.8, 92.1,
  (SELECT id FROM object_types WHERE name = 'stair'), '4ccb7c2d-0bb2-49e0-828f-5b2c43787ee5');

-- 2 этаж — Лестница 2 (центральная левая)
INSERT INTO objects (id, name, alias, x, y, width, height, object_type_id, floor_id)
VALUES ('52454b36-1895-41b4-b39b-d8f71a5ff9be', 'Лестница 2', 'Центральная левая', 753.3, 286.7, 92.1, 67.3,
  (SELECT id FROM object_types WHERE name = 'stair'), '4ccb7c2d-0bb2-49e0-828f-5b2c43787ee5');

-- 2 этаж — Лестница 3 (центральная правая)
INSERT INTO objects (id, name, alias, x, y, width, height, object_type_id, floor_id)
VALUES ('a41ffeb5-d54c-47e3-88a5-2ad3059fe957', 'Лестница 3', 'Центральная правая', 964.8, 286.7, 92.1, 67.3,
  (SELECT id FROM object_types WHERE name = 'stair'), '4ccb7c2d-0bb2-49e0-828f-5b2c43787ee5');

-- 2 этаж — Лестница 4 (северо-восточная)
INSERT INTO objects (id, name, alias, x, y, width, height, object_type_id, floor_id)
VALUES ('2dcd6de3-e4a4-434f-b9b3-e85c3d2ab937', 'Лестница 4', 'Северо-восточная секция', 864.8, 1101.9, 41.6, 95.9,
  (SELECT id FROM object_types WHERE name = 'stair'), '4ccb7c2d-0bb2-49e0-828f-5b2c43787ee5');

-- 2 этаж — Лестница 5 (правый верхний угол)
INSERT INTO objects (id, name, alias, x, y, width, height, object_type_id, floor_id)
VALUES ('5290b3b6-f24f-4320-9569-69b68a750536', 'Лестница 5', 'Правый верхний угол', 1816.4, 903, 87.6, 107.5,
  (SELECT id FROM object_types WHERE name = 'stair'), '4ccb7c2d-0bb2-49e0-828f-5b2c43787ee5');

-- 2 этаж — Лестница 6 (входная)
INSERT INTO objects (id, name, alias, x, y, width, height, object_type_id, floor_id)
VALUES ('d92d4131-07cf-4534-8419-2de0f91ebf3e', 'Лестница 6', 'Входная лестница', 1589.5, 56.5, 67.3, 92.1,
  (SELECT id FROM object_types WHERE name = 'stair'), '4ccb7c2d-0bb2-49e0-828f-5b2c43787ee5');

-- 3 этаж — Лестница 1 (западная)
INSERT INTO objects (id, name, alias, x, y, width, height, object_type_id, floor_id)
VALUES ('8772b21a-e183-4fa5-a3ed-2e011aebef84', 'Лестница 1', 'Западная лестница', 247.1, 1074.5, 66.8, 92.1,
  (SELECT id FROM object_types WHERE name = 'stair'), 'acc38768-c209-4702-957d-778e59f875f3');

-- 3 этаж — Лестница 2 (центральная левая)
INSERT INTO objects (id, name, alias, x, y, width, height, object_type_id, floor_id)
VALUES ('1f1c888e-3dc4-4f86-90e6-41e27ee7045e', 'Лестница 2', 'Центральная левая', 753.3, 286.8, 92.1, 67.3,
  (SELECT id FROM object_types WHERE name = 'stair'), 'acc38768-c209-4702-957d-778e59f875f3');

-- 3 этаж — Лестница 3 (центральная правая)
INSERT INTO objects (id, name, alias, x, y, width, height, object_type_id, floor_id)
VALUES ('79a4a4d2-e991-4b4a-9f58-facc68c45b33', 'Лестница 3', 'Центральная правая', 964.8, 286.8, 92.1, 67.3,
  (SELECT id FROM object_types WHERE name = 'stair'), 'acc38768-c209-4702-957d-778e59f875f3');

-- 4 этаж — Лестница 1 (у гардероба)
INSERT INTO objects (id, name, alias, x, y, width, height, object_type_id, floor_id)
VALUES ('775f5c3c-d5a5-4e01-b933-dc4236e4d423', 'Лестница у гардероба', 'Лестница у гардероба', 247.1, 1074.5, 66.8, 92.1,
  (SELECT id FROM object_types WHERE name = 'stair'), '517650f2-3e3a-42d0-b6e9-6c8e7271b096');

-- 4 этаж — Лестница 2 (у туалета)
INSERT INTO objects (id, name, alias, x, y, width, height, object_type_id, floor_id)
VALUES ('67b40b8e-64b1-4c80-b4a2-857e3b7fa7be', 'Лестница у туалета', 'Лестница у туалета', 753.2, 286.8, 92.1, 67.3,
  (SELECT id FROM object_types WHERE name = 'stair'), '517650f2-3e3a-42d0-b6e9-6c8e7271b096');


---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

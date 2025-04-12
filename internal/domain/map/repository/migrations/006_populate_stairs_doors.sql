-- Двери 1 этажа для лестниц.
INSERT INTO doors (id, x, y, width, height) VALUES ('5b15bb73-178a-48cd-8332-84939e0e5fd1', 259.6, 1073.4, 41.6, 3);
INSERT INTO doors (id, x, y, width, height) VALUES ('48679011-7ee6-49be-9dc1-451251f2838f', 455.7, 987.1, 15, 3);
INSERT INTO doors (id, x, y, width, height) VALUES ('c8e98c83-3556-47bf-bf95-005d8f2d1793', 842.9, 849.2, 3, 20);
INSERT INTO doors (id, x, y, width, height) VALUES ('a37e6006-afa9-4c8f-923f-e386d128185d', 843.9, 407.1, 3, 20);
INSERT INTO doors (id, x, y, width, height) VALUES ('6a6c2c20-3965-49fd-9bc4-f768695d0a54', 844, 304.2, 3, 33.6);
INSERT INTO doors (id, x, y, width, height) VALUES ('b606832a-5a0a-458c-b6ee-c33770be3a02', 963.3, 304.2, 3, 33.6);
INSERT INTO doors (id, x, y, width, height) VALUES ('e548e803-6db1-47b0-9edb-45efa0a15a5e', 875.5, 1105.3, 20, 3);
INSERT INTO doors (id, x, y, width, height) VALUES ('9c070c69-0338-411a-b9be-3f5fbdf9f4d1', 1372.7, 784.5, 20, 3);
INSERT INTO doors (id, x, y, width, height) VALUES ('0b2a19e0-4334-43c3-9c29-0945055b9abe', 1814.9, 946.9, 3, 20);
INSERT INTO doors (id, x, y, width, height) VALUES ('ef127f51-5331-498f-a5c7-c4b8fc4faa7f', 1606.1, 147.7, 33.6, 3);

INSERT INTO object_doors (object_id, door_id) VALUES ('a98ef279-e4a7-44e4-b093-71221e3b8e4f', '5b15bb73-178a-48cd-8332-84939e0e5fd1');
INSERT INTO object_doors (object_id, door_id) VALUES ('9b728658-f5fd-4f31-b872-90f60e1eac36', '48679011-7ee6-49be-9dc1-451251f2838f');
INSERT INTO object_doors (object_id, door_id) VALUES ('fd8d2a93-55e3-4a5f-87a0-c9a12f0021cb', 'c8e98c83-3556-47bf-bf95-005d8f2d1793');
INSERT INTO object_doors (object_id, door_id) VALUES ('1ed2ed7d-9643-4b7f-a982-f296c8cccffe', 'a37e6006-afa9-4c8f-923f-e386d128185d');
INSERT INTO object_doors (object_id, door_id) VALUES ('b71ae235-8b14-4175-9607-3068d49a3c8a', '6a6c2c20-3965-49fd-9bc4-f768695d0a54');
INSERT INTO object_doors (object_id, door_id) VALUES ('87d3a8f5-2795-4530-998e-20bbbc48240e', 'b606832a-5a0a-458c-b6ee-c33770be3a02');
INSERT INTO object_doors (object_id, door_id) VALUES ('20dcdfbd-2cd7-4861-bb0c-3be1ab1504a1', 'e548e803-6db1-47b0-9edb-45efa0a15a5e');
INSERT INTO object_doors (object_id, door_id) VALUES ('8b4c9b42-80fd-4a85-9d52-ef78a948f2e4', '9c070c69-0338-411a-b9be-3f5fbdf9f4d1');
INSERT INTO object_doors (object_id, door_id) VALUES ('2f04a0ae-69b3-4534-9061-97c340e261f3', '0b2a19e0-4334-43c3-9c29-0945055b9abe');
INSERT INTO object_doors (object_id, door_id) VALUES ('82a1a2b2-356e-4f61-a07d-3e6f48a6b2a2', 'ef127f51-5331-498f-a5c7-c4b8fc4faa7f');

-- Двери 2 этажа для лестниц.
INSERT INTO doors (id, x, y, width, height) VALUES ('6d53392d-ce17-46b0-9f47-7bf2bfbea3be', 270.6, 1072.5, 20, 3);
INSERT INTO doors (id, x, y, width, height) VALUES ('b7a2f238-df66-4a7a-877a-ef734a2f234a', 843.7, 310.4, 3, 20);
INSERT INTO doors (id, x, y, width, height) VALUES ('afb5fde7-f216-43b6-9803-9e154a223904', 963.3, 310.3, 3, 20);
INSERT INTO doors (id, x, y, width, height) VALUES ('ecc88e98-d11c-4669-857b-c4cf6f6d2c9a', 1614.7, 147.1, 20, 3);
INSERT INTO doors (id, x, y, width, height) VALUES ('bddf2484-fc4b-40c7-8af9-f51c7a363c6f', 1814.9, 948.5, 3, 20);
INSERT INTO doors (id, x, y, width, height) VALUES ('fb32c57c-95a1-439b-99b0-dc75371178da', 875.6, 1196.3, 20, 3);
INSERT INTO object_doors (object_id, door_id) VALUES ('2ae8c611-2321-46d6-933f-bb5e17c6bb38', '6d53392d-ce17-46b0-9f47-7bf2bfbea3be');
INSERT INTO object_doors (object_id, door_id) VALUES ('52454b36-1895-41b4-b39b-d8f71a5ff9be', 'b7a2f238-df66-4a7a-877a-ef734a2f234a');
INSERT INTO object_doors (object_id, door_id) VALUES ('a41ffeb5-d54c-47e3-88a5-2ad3059fe957', 'afb5fde7-f216-43b6-9803-9e154a223904');
INSERT INTO object_doors (object_id, door_id) VALUES ('d92d4131-07cf-4534-8419-2de0f91ebf3e', 'ecc88e98-d11c-4669-857b-c4cf6f6d2c9a');
INSERT INTO object_doors (object_id, door_id) VALUES ('5290b3b6-f24f-4320-9569-69b68a750536', 'bddf2484-fc4b-40c7-8af9-f51c7a363c6f');
INSERT INTO object_doors (object_id, door_id) VALUES ('2dcd6de3-e4a4-434f-b9b3-e85c3d2ab937', 'fb32c57c-95a1-439b-99b0-dc75371178da');

-- Двери 3 этажа для лестниц.
INSERT INTO doors (id, x, y, width, height) VALUES ('867bdbb4-0b55-46e3-b416-8f75080f4e10', 266.8, 1073, 20, 3);
INSERT INTO doors (id, x, y, width, height) VALUES ('af9ec40c-85f8-4d86-adfb-c67c7b8cee57', 843.9, 302.4, 3, 20);
INSERT INTO doors (id, x, y, width, height) VALUES ('957f21d9-ced2-4f40-b302-99f26995b605', 963.3, 302.4, 3, 20);
INSERT INTO object_doors (object_id, door_id) VALUES ('8772b21a-e183-4fa5-a3ed-2e011aebef84', '867bdbb4-0b55-46e3-b416-8f75080f4e10');
INSERT INTO object_doors (object_id, door_id) VALUES ('1f1c888e-3dc4-4f86-90e6-41e27ee7045e', 'af9ec40c-85f8-4d86-adfb-c67c7b8cee57');
INSERT INTO object_doors (object_id, door_id) VALUES ('79a4a4d2-e991-4b4a-9f58-facc68c45b33', '957f21d9-ced2-4f40-b302-99f26995b605');

-- Двери 4 этажа для лестниц.
INSERT INTO doors (id, x, y, width, height) VALUES ('72466f92-2d62-465e-a89c-e5020a43c168', 268.3, 1073, 20, 3);
INSERT INTO doors (id, x, y, width, height) VALUES ('68fb988f-3ff1-46e6-9ef1-561c14755614', 843.8, 303.6, 3, 20);
INSERT INTO object_doors (object_id, door_id) VALUES ('775f5c3c-d5a5-4e01-b933-dc4236e4d423', '72466f92-2d62-465e-a89c-e5020a43c168');
INSERT INTO object_doors (object_id, door_id) VALUES ('67b40b8e-64b1-4c80-b4a2-857e3b7fa7be', '68fb988f-3ff1-46e6-9ef1-561c14755614');

---- create above / drop below ----

DELETE FROM object_doors WHERE door_id IN (
    '72466f92-2d62-465e-a89c-e5020a43c168',
    '68fb988f-3ff1-46e6-9ef1-561c14755614',
    '867bdbb4-0b55-46e3-b416-8f75080f4e10',
    'af9ec40c-85f8-4d86-adfb-c67c7b8cee57',
    '957f21d9-ced2-4f40-b302-99f26995b605',
    '6d53392d-ce17-46b0-9f47-7bf2bfbea3be',
    'b7a2f238-df66-4a7a-877a-ef734a2f234a',
    'afb5fde7-f216-43b6-9803-9e154a223904',
    'ecc88e98-d11c-4669-857b-c4cf6f6d2c9a',
    'bddf2484-fc4b-40c7-8af9-f51c7a363c6f',
    'fb32c57c-95a1-439b-99b0-dc75371178da',
    '5b15bb73-178a-48cd-8332-84939e0e5fd1',
    '48679011-7ee6-49be-9dc1-451251f2838f',
    'c8e98c83-3556-47bf-bf95-005d8f2d1793',
    'a37e6006-afa9-4c8f-923f-e386d128185d',
    '6a6c2c20-3965-49fd-9bc4-f768695d0a54',
    'b606832a-5a0a-458c-b6ee-c33770be3a02',
    'e548e803-6db1-47b0-9edb-45efa0a15a5e',
    '9c070c69-0338-411a-b9be-3f5fbdf9f4d1',
    '0b2a19e0-4334-43c3-9c29-0945055b9abe',
    'ef127f51-5331-498f-a5c7-c4b8fc4faa7f'
);

DELETE FROM doors WHERE id IN (
    '72466f92-2d62-465e-a89c-e5020a43c168',
    '68fb988f-3ff1-46e6-9ef1-561c14755614',
    '867bdbb4-0b55-46e3-b416-8f75080f4e10',
    'af9ec40c-85f8-4d86-adfb-c67c7b8cee57',
    '957f21d9-ced2-4f40-b302-99f26995b605',
    '6d53392d-ce17-46b0-9f47-7bf2bfbea3be',
    'b7a2f238-df66-4a7a-877a-ef734a2f234a',
    'afb5fde7-f216-43b6-9803-9e154a223904',
    'ecc88e98-d11c-4669-857b-c4cf6f6d2c9a',
    'bddf2484-fc4b-40c7-8af9-f51c7a363c6f',
    'fb32c57c-95a1-439b-99b0-dc75371178da',
    '5b15bb73-178a-48cd-8332-84939e0e5fd1',
    '48679011-7ee6-49be-9dc1-451251f2838f',
    'c8e98c83-3556-47bf-bf95-005d8f2d1793',
    'a37e6006-afa9-4c8f-923f-e386d128185d',
    '6a6c2c20-3965-49fd-9bc4-f768695d0a54',
    'b606832a-5a0a-458c-b6ee-c33770be3a02',
    'e548e803-6db1-47b0-9edb-45efa0a15a5e',
    '9c070c69-0338-411a-b9be-3f5fbdf9f4d1',
    '0b2a19e0-4334-43c3-9c29-0945055b9abe',
    'ef127f51-5331-498f-a5c7-c4b8fc4faa7f'
);

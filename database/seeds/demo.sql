INSERT INTO objects(guid, type_id, name, created_at, updated_at)
VALUES
  ('c949de70-7e44-42b7-b34e-b84efd1afbf1', 1, 'Button', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('c949de70-7e44-42b7-b34e-b84efd1afbf2', 2, 'Display', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('c949de70-7e44-42b7-b34e-b84efd1afbf3', 3, 'ComplexObject', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('c949de70-7e44-42b7-b34e-b84efd1afbf5', 4, 'Valve', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('c949de70-7e44-42b7-b34e-b84efd1afbf6', 5, 'Zone', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('c949de70-7e44-42b7-b34e-b84efd1afbf7', 6, 'Animated', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('c949de70-7e44-42b7-b34e-b84efd1afbf8', 7, 'Liquid', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO worlds(name, created_at, updated_at)
VALUES
  ('Demo world', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO locations(guid, name, created_at, updated_at)
VALUES
  ('a949de70-7e44-42b7-b34e-b84efd1afbf1', 'Demo location', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO world_locations(name, world_id, location_id, created_at, updated_at)
VALUES
   ( 'Demo location', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO groups(name, world_location_id, created_at, updated_at)
VALUES
  ('Common', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
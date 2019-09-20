CREATE TABLE world_locations (
  id bigserial PRIMARY KEY,
  sid uuid NOT NULL,
  created_at timestamp NOT NULL,
  updated_at timestamp NOT NULL,
  world_id bigint NOT NULL,
  location_id bigint NOT NULL,
  name character varying(255) NOT NULL
);

CREATE UNIQUE INDEX idx_world_locations_sid ON world_locations USING btree (world_id, sid);

ALTER TABLE world_locations
  ADD CONSTRAINT fk_world_locations_to_worlds
FOREIGN KEY (world_id) REFERENCES worlds(id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE world_locations
  ADD CONSTRAINT fk_world_locations_to_locations
FOREIGN KEY (location_id) REFERENCES locations(id) ON UPDATE CASCADE ON DELETE CASCADE;
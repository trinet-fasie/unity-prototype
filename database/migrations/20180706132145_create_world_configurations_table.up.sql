CREATE TABLE world_configurations (
    id bigserial PRIMARY KEY,
    sid uuid NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    world_id bigint NOT NULL,
    start_world_location_id bigint NOT NULL,
    name character varying(255)
);

CREATE UNIQUE INDEX idx_world_configurations_sid ON world_configurations USING btree (world_id, sid);

ALTER TABLE world_configurations
ADD CONSTRAINT fk_world_configurations_to_worlds
FOREIGN KEY (world_id) REFERENCES worlds(id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE world_configurations
ADD CONSTRAINT fk_world_configurations_to_world_locations
FOREIGN KEY (start_world_location_id) REFERENCES world_locations(id) ON UPDATE CASCADE ON DELETE RESTRICT;

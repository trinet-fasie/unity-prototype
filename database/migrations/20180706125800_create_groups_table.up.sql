CREATE TABLE groups (
    id bigserial PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    world_location_id bigint NOT NULL,
    code text DEFAULT '',
    editor_data json NOT NULL DEFAULT '{}',
    name character varying(255)
);

ALTER TABLE groups
ADD CONSTRAINT fk_groups_to_locations
FOREIGN KEY (world_location_id) REFERENCES world_locations(id) ON UPDATE CASCADE ON DELETE CASCADE;

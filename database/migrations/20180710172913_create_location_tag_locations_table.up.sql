CREATE TABLE location_tag_locations (
  location_tag_id bigserial NOT NULL,
  location_id bigserial NOT NULL,
  created_at timestamp NOT NULL
);

ALTER TABLE ONLY location_tag_locations
  ADD CONSTRAINT idx_location_tag_locations PRIMARY KEY (location_tag_id, location_id);

ALTER TABLE location_tag_locations
  ADD CONSTRAINT fk_location_tag_locations_to_location_tag
FOREIGN KEY (location_tag_id) REFERENCES location_tags(id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE location_tag_locations
  ADD CONSTRAINT fk_location_tag_locations_to_location
FOREIGN KEY (location_id) REFERENCES locations(id) ON UPDATE CASCADE ON DELETE RESTRICT;

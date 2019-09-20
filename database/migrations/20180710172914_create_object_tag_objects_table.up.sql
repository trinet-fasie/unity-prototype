CREATE TABLE object_tag_objects (
  object_tag_id bigserial NOT NULL,
  object_id bigserial NOT NULL,
  created_at timestamp NOT NULL
);

ALTER TABLE ONLY object_tag_objects
  ADD CONSTRAINT idx_object_tag_objects PRIMARY KEY (object_tag_id, object_id);

ALTER TABLE object_tag_objects
  ADD CONSTRAINT fk_object_tag_objects_to_object_tag
FOREIGN KEY (object_tag_id) REFERENCES object_tags(id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE object_tag_objects
  ADD CONSTRAINT fk_object_tag_objects_to_object
FOREIGN KEY (object_id) REFERENCES objects(id) ON UPDATE CASCADE ON DELETE RESTRICT;

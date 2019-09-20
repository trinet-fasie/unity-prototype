CREATE TABLE group_objects (
    id bigserial PRIMARY KEY,
    parent_id bigint DEFAULT NULL,
    position int DEFAULT 0 NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    group_id bigint NOT NULL,
    object_id bigint NOT NULL,
    instance_id bigint NOT NULL,
    name character varying(255) NOT NULL,
    data json NOT NULL,
    locked boolean NOT NULL DEFAULT FALSE
);

ALTER TABLE group_objects
ADD CONSTRAINT fk_group_objects_to_groups
FOREIGN KEY (group_id) REFERENCES groups(id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE group_objects
ADD CONSTRAINT fk_group_objects_to_group_object
FOREIGN KEY (parent_id) REFERENCES group_objects(id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE group_objects
ADD CONSTRAINT fk_group_objects_to_objects
FOREIGN KEY (object_id) REFERENCES objects(id) ON UPDATE CASCADE ON DELETE CASCADE;


CREATE UNIQUE INDEX idx_group_objects_unique_name ON group_objects USING btree (group_id, name);

CREATE UNIQUE INDEX idx_group_objects_unique_instance_id ON group_objects USING btree (group_id, instance_id);

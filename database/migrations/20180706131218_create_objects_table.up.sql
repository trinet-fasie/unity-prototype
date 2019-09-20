CREATE TABLE objects (
    id bigserial PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    guid uuid NOT NULL,
    config json NOT NULL DEFAULT '{}'
);

CREATE UNIQUE INDEX idx_object_guid ON objects USING btree (guid);
CREATE TABLE location_tags (
    id bigserial PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    text character varying(255) NOT NULL
);

CREATE UNIQUE INDEX idx_location_tag_text ON location_tags USING btree (text);
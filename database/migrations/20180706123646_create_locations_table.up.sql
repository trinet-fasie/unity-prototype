CREATE TABLE locations (
    id bigserial PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    guid uuid NOT NULL,
    name character varying(255)
);
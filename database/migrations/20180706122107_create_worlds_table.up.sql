CREATE TABLE worlds (
    id bigserial PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    name character varying(255) NOT NULL
);
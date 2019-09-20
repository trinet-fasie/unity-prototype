CREATE TABLE world_configuration_groups (
    world_configuration_id bigint NOT NULL,
    group_id bigint NOT NULL,
    created_at timestamp NOT NULL
);

ALTER TABLE ONLY world_configuration_groups
ADD CONSTRAINT idx_world_configuration_groups PRIMARY KEY (world_configuration_id, group_id);

ALTER TABLE world_configuration_groups
ADD CONSTRAINT fk_world_configuration_groups_to_world_configuration
FOREIGN KEY (world_configuration_id) REFERENCES world_configurations(id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE world_configuration_groups
ADD CONSTRAINT fk_world_configuration_groups_to_groups
FOREIGN KEY (group_id) REFERENCES groups(id) ON UPDATE CASCADE ON DELETE RESTRICT;

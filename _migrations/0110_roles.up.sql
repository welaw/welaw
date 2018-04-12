CREATE TABLE roles (
    uid             UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    name            varchar(255) NOT NULL,
    created_at      TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP WITHOUT TIME ZONE
);

CREATE UNIQUE INDEX roles_name_idx ON roles(LOWER(name)) WHERE deleted_at IS NULL;

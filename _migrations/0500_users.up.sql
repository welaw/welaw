CREATE TABLE users (
    key                 UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    uid                 UUID NOT NULL DEFAULT gen_random_uuid(),
    provider            UUID NOT NULL,
    provider_id         varchar(255) NOT NULL DEFAULT '',
    username            varchar(255),
    full_name           varchar(255) NOT NULL,
    full_name_private   boolean NOT NULL DEFAULT true,
    email               varchar(255) NOT NULL,
    email_private       boolean NOT NULL DEFAULT true,
    picture_url         varchar(255) NOT NULL,
    biography           varchar NOT NULL,
    url                 varchar(255) NOT NULL DEFAULT '',
    upstream            UUID,
    password            varchar,
    metadata            jsonb NOT NULL DEFAULT '{}'::jsonb,

    groups              TEXT NOT NULL DEFAULT '',

    last_login          TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at          TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at          TIMESTAMP WITHOUT TIME ZONE
);

-- CREATE UNIQUE INDEX users_email_idx ON users(LOWER(email)) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX users_provider_id_idx ON users(provider_id) WHERE deleted_at IS NULL AND provider_id <> '';
CREATE UNIQUE INDEX users_uid ON users(uid) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX users_username_idx ON users(LOWER(username)) WHERE deleted_at IS NULL;

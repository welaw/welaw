CREATE TABLE upstreams (
    uid                     UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    ident                   varchar(255) NOT NULL,
    upstream_name           varchar(255) NOT NULL,
    upstream_description    TEXT NOT NULL,
    upstream_url            varchar(255) NOT NULL DEFAULT '',
    name                    varchar(255) NOT NULL,
    description             TEXT NOT NULL,
    url                     varchar(255) NOT NULL,
    default_user            UUID NOT NULL DEFAULT uuid_nil(),
    geo_coords              POINT,
    metadata                JSONB NOT NULL DEFAULT '{}'::jsonb,
    tags                    JSONB NOT NULL DEFAULT '{}'::jsonb,
    user_id                 UUID NOT NULL,
    created_at              TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at              TIMESTAMP WITHOUT TIME ZONE
);


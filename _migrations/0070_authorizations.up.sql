CREATE TABLE authorizations (
    uid                 UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id             UUID NOT NULL,
    provider            UUID NOT NULL,
    provider_id         varchar(255) NOT NULL DEFAULT '',
    password            varchar,
    created_at          TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at          TIMESTAMP WITHOUT TIME ZONE
);

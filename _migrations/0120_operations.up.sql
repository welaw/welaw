CREATE TABLE operations (
    uid                  UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    name                 varchar(255) NOT NULL,
    description          varchar(255) NOT NULL,
    created_at           TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at           TIMESTAMP WITHOUT TIME ZONE
);

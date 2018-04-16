CREATE TABLE role_operations (
    role_id              UUID NOT NULL,
    operation_id         UUID NOT NULL,
    objects              varchar(255),
    created_at           TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at           TIMESTAMP WITHOUT TIME ZONE
);

CREATE TABLE user_roles (
    user_id              UUID NOT NULL,
    role_id              UUID NOT NULL,
    created_at           TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at           TIMESTAMP WITHOUT TIME ZONE
);

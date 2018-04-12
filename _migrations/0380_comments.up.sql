CREATE TABLE comments (
  key                       UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  uid                       UUID NOT NULL DEFAULT gen_random_uuid(),
  version_id                UUID NOT NULL,
  user_id                   UUID NOT NULL,
  comment                   TEXT NOT NULL DEFAULT '',
  disabled                  BOOLEAN NOT NULL DEFAULT TRUE,
  created_at                TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at                TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at                TIMESTAMP WITHOUT TIME ZONE
);

CREATE UNIQUE INDEX comments_version_user_idx ON comments(version_id, user_id) WHERE deleted_at IS NULL;

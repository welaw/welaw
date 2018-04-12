CREATE TYPE vote_value AS ENUM('YES', 'NO', 'PRESENT', 'NOT_PRESENT');

CREATE TABLE votes (
  uid           UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  -- user_id       UUID NOT NULL REFERENCES users(uid) ON DELETE CASCADE,
  -- version_id    UUID NOT NULL REFERENCES versions(uid) ON DELETE CASCADE,
  user_id       UUID NOT NULL,
  version_id    UUID NOT NULL,
  value         vote_value DEFAULT NULL,
  -- vote          varchar NOT NULL,
  comment       TEXT NOT NULL,
  cast_at       TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at    TIMESTAMP WITHOUT TIME ZONE
);

CREATE UNIQUE INDEX ON votes (user_id, version_id) WHERE deleted_at IS NULL;

CREATE TABLE annotations (
  uid                       UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  comment_id                UUID NOT NULL,
  id                        UUID NOT NULL DEFAULT gen_random_uuid(),
  text                      VARCHAR(255) NOT NULL,
  quote                     TEXT NOT NULL,
  ranges                    JSONB NOT NULL,
  comment                   TEXT NOT NULL DEFAULT '',
  created_at                TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at                TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at                TIMESTAMP WITHOUT TIME ZONE
);

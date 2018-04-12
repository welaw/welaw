CREATE TABLE upstream_tags (
  uid               UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  upstream_id       UUID NOT NULL,
  ident             VARCHAR(255) NOT NULL,
  ranking           INTEGER NOT NULL,
  name              VARCHAR(255) NOT NULL,
  description       TEXT NOT NULL,
  number_type       BOOLEAN NOT NULL,
  created_at        TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at        TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at        TIMESTAMP WITHOUT TIME ZONE
);

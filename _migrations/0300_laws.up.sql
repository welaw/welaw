CREATE TABLE laws (
  uid                       UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  upstream_id               UUID NOT NULL,
  user_id                   UUID NOT NULL,
  ident                     varchar(255) NOT NULL,
  title                     TEXT NOT NULL,
  short_title               varchar(255) NOT NULL,
  description               TEXT NOT NULL,
  url                       varchar(255),
  upstream_group_id         UUID NOT NULL,
  published_at              TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_at                TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at                TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at                TIMESTAMP WITHOUT TIME ZONE
);

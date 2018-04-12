CREATE TABLE versions (
  uid                   UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  branch_id             UUID NOT NULL,
  user_id               UUID NOT NULL,
  upstream_group_id     UUID,
  hash                  varchar(255) NOT NULL,
  message               varchar(255) NOT NULL,
  number                integer NOT NULL DEFAULT 1,
  tag_1                 varchar(255),
  tag_2                 varchar(255),
  tag_3                 varchar(255),
  tag_4                 varchar(255),
  published_at          TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_at            TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at            TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at            TIMESTAMP WITHOUT TIME ZONE
);

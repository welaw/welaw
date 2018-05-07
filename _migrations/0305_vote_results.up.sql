CREATE TYPE vote_result AS ENUM('PASSED', 'FAILED', 'SIGNED');

CREATE TABLE vote_results (
  uid                       UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  user_id                   UUID NOT NULL,
  law_id                    UUID NOT NULL,
  result                    vote_result NOT NULL,
  published_at              TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_at                TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at                TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at                TIMESTAMP WITHOUT TIME ZONE
);

CREATE TABLE comment_likes (
  uid                       UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  comment_id                UUID NOT NULL,
  user_id                   UUID NOT NULL,
  created_at                TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at                TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at                TIMESTAMP WITHOUT TIME ZONE
);

CREATE UNIQUE INDEX comment_likes_comment_user_idx ON comment_likes(comment_id, user_id) WHERE deleted_at IS NULL;

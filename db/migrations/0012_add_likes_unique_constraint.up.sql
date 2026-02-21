ALTER TABLE likes ADD CONSTRAINT uq_likes_sake_token UNIQUE (sake_id, token);

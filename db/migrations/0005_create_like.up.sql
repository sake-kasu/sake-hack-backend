CREATE TABLE likes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sake_id UUID NOT NULL REFERENCES sake(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_likes_sake_id ON likes(sake_id);
CREATE INDEX idx_likes_token ON likes(token);

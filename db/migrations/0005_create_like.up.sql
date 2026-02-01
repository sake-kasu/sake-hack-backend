CREATE TABLE likes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sake_id UUID NOT NULL REFERENCES sake(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_likes_sake_id ON likes(sake_id);
CREATE INDEX idx_likes_token ON likes(token);

COMMENT ON TABLE likes IS '酒へのいいねを管理するテーブル';

COMMENT ON COLUMN likes.id IS 'いいねの一意識別子';
COMMENT ON COLUMN likes.sake_id IS '酒ID';
COMMENT ON COLUMN likes.token IS 'ユーザーを識別するトークン';
COMMENT ON COLUMN likes.created_at IS '作成日時';

CREATE TABLE bookmarks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sake_id UUID NOT NULL REFERENCES sakes(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(user_id, sake_id)
);

CREATE INDEX idx_bookmarks_user_id ON bookmarks(user_id);
CREATE INDEX idx_bookmarks_sake_id ON bookmarks(sake_id);

COMMENT ON TABLE bookmarks IS 'ユーザーのブックマークを管理するテーブル';

COMMENT ON COLUMN bookmarks.id IS 'ブックマークの一意識別子';
COMMENT ON COLUMN bookmarks.user_id IS 'ユーザーID';
COMMENT ON COLUMN bookmarks.sake_id IS '酒ID';
COMMENT ON COLUMN bookmarks.created_at IS '作成日時';

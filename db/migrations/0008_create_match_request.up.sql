CREATE TYPE match_status AS ENUM ('PENDING', 'ACCEPTED', 'REJECTED');

CREATE TABLE match_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    to_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status match_status NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_match_requests_from_user_id ON match_requests(from_user_id);
CREATE INDEX idx_match_requests_to_user_id ON match_requests(to_user_id);
CREATE INDEX idx_match_requests_status ON match_requests(status);

COMMENT ON TYPE match_status IS 'マッチングリクエストのステータスを定義するENUM';

COMMENT ON TABLE match_requests IS 'ユーザー間のマッチングリクエストを管理するテーブル';

COMMENT ON COLUMN match_requests.id IS 'マッチングリクエストの一意識別子';
COMMENT ON COLUMN match_requests.from_user_id IS 'リクエスト送信者ID';
COMMENT ON COLUMN match_requests.to_user_id IS 'リクエスト受信者ID';
COMMENT ON COLUMN match_requests.status IS 'マッチングリクエストのステータス(PENDING: 保留中, ACCEPTED: 承認済み, REJECTED: 拒否)';
COMMENT ON COLUMN match_requests.created_at IS '作成日時';
COMMENT ON COLUMN match_requests.updated_at IS '更新日時';

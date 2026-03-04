CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE sake_category AS ENUM (
    'JAPANESE_SAKE',
    'SHOCHU',
    'AWAMORI',
    'BEER',
    'WINE',
    'FRUIT_WINE',
    'LIQUEUR',
    'NON_ALCOHOL',
    'OTHER'
);

CREATE TYPE user_role AS ENUM ('SHOP', 'USER');

CREATE TYPE match_status AS ENUM ('PENDING', 'ACCEPTED', 'REJECTED');

CREATE TABLE sakes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL,
    phonetic VARCHAR(100),
    category sake_category NOT NULL,
    alcohol_percentage NUMERIC(4, 1) CHECK (alcohol_percentage >= 0 AND alcohol_percentage <= 100),
    volume_max INTEGER CHECK (volume_max > 0),
    volume_remain INTEGER CHECK (volume_remain >= 0 AND volume_remain <= 100),
    region VARCHAR(100),
    price INTEGER CHECK (price >= 0),
    memo VARCHAR(500),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_sakes_category ON sakes(category);
CREATE INDEX idx_sakes_name ON sakes(name);

CREATE TABLE sake_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tag VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE sake_tag_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sake_id UUID NOT NULL REFERENCES sakes(id) ON DELETE CASCADE,
    sake_tag_id UUID NOT NULL REFERENCES sake_tags(id) ON DELETE CASCADE,
    CONSTRAINT uq_sake_tag_links_sake_tag UNIQUE (sake_id, sake_tag_id)
);
CREATE INDEX idx_sake_tag_links_sake_id ON sake_tag_links(sake_id);
CREATE INDEX idx_sake_tag_links_sake_tag_id ON sake_tag_links(sake_tag_id);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(100),
    role user_role NOT NULL DEFAULT 'USER',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_users_email ON users(email);

CREATE TABLE likes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sake_id UUID NOT NULL REFERENCES sakes(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_likes_sake_token UNIQUE (sake_id, token)
);
CREATE INDEX idx_likes_sake_id ON likes(sake_id);
CREATE INDEX idx_likes_token ON likes(token);

CREATE TABLE drink_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sake_id UUID NOT NULL REFERENCES sakes(id) ON DELETE CASCADE,
    drank_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_drink_logs_user_id ON drink_logs(user_id);
CREATE INDEX idx_drink_logs_sake_id ON drink_logs(sake_id);
CREATE INDEX idx_drink_logs_drank_at ON drink_logs(drank_at);

CREATE TABLE bookmarks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sake_id UUID NOT NULL REFERENCES sakes(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_bookmarks_user_sake UNIQUE (user_id, sake_id)
);
CREATE INDEX idx_bookmarks_user_id ON bookmarks(user_id);
CREATE INDEX idx_bookmarks_sake_id ON bookmarks(sake_id);

CREATE TABLE match_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    to_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status match_status NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_match_requests_not_self CHECK (from_user_id <> to_user_id)
);
CREATE INDEX idx_match_requests_from_user_id ON match_requests(from_user_id);
CREATE INDEX idx_match_requests_to_user_id ON match_requests(to_user_id);
CREATE INDEX idx_match_requests_status ON match_requests(status);

CREATE TABLE sake_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sake_id UUID NOT NULL REFERENCES sakes(id) ON DELETE CASCADE,
    image_key VARCHAR(500) NOT NULL,
    sort_order INTEGER NOT NULL CHECK (sort_order >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_sake_images_sake_sort_order UNIQUE (sake_id, sort_order),
    CONSTRAINT uq_sake_images_sake_image_key UNIQUE (sake_id, image_key)
);
CREATE INDEX idx_sake_images_sake_id ON sake_images(sake_id);
CREATE INDEX idx_sake_images_sake_sort_order ON sake_images(sake_id, sort_order);

COMMENT ON TYPE sake_category IS 'お酒カテゴリ: JAPANESE_SAKE/SHOCHU/AWAMORI/BEER/WINE/FRUIT_WINE/LIQUEUR/NON_ALCOHOL/OTHER';
COMMENT ON TYPE user_role IS 'ユーザーロール: SHOP(主催者)/USER(参加者)';
COMMENT ON TYPE match_status IS 'マッチングリクエスト状態: PENDING/ACCEPTED/REJECTED';

COMMENT ON TABLE sakes IS 'お酒マスタ';
COMMENT ON COLUMN sakes.id IS '酒ID';
COMMENT ON COLUMN sakes.name IS '酒名（最大50文字, 必須）';
COMMENT ON COLUMN sakes.phonetic IS 'ふりがな（最大100文字）';
COMMENT ON COLUMN sakes.category IS '酒カテゴリ（必須）';
COMMENT ON COLUMN sakes.alcohol_percentage IS 'アルコール度数（0.0-100.0）';
COMMENT ON COLUMN sakes.volume_max IS '最大容量（ml）';
COMMENT ON COLUMN sakes.volume_remain IS '残り容量（%: 0-100）';
COMMENT ON COLUMN sakes.region IS '産地';
COMMENT ON COLUMN sakes.price IS '価格（円）';
COMMENT ON COLUMN sakes.memo IS 'メモ（最大500文字）';
COMMENT ON COLUMN sakes.created_at IS '作成日時';
COMMENT ON COLUMN sakes.updated_at IS '更新日時';

COMMENT ON TABLE sake_tags IS 'お酒タグマスタ';
COMMENT ON COLUMN sake_tags.id IS 'タグID';
COMMENT ON COLUMN sake_tags.tag IS 'タグ文字列（最大50文字, 必須）';

COMMENT ON TABLE sake_tag_links IS 'お酒とタグの関連';
COMMENT ON COLUMN sake_tag_links.id IS '一意識別子';
COMMENT ON COLUMN sake_tag_links.sake_id IS 'お酒ID';
COMMENT ON COLUMN sake_tag_links.sake_tag_id IS 'タグID';

COMMENT ON TABLE users IS 'ユーザー';
COMMENT ON COLUMN users.id IS '一意識別子';
COMMENT ON COLUMN users.email IS 'メールアドレス（必須, 一意）';
COMMENT ON COLUMN users.password_hash IS 'パスワードハッシュ（必須）';
COMMENT ON COLUMN users.display_name IS '表示名（任意）';
COMMENT ON COLUMN users.role IS 'ロール（SHOP/USER）';
COMMENT ON COLUMN users.created_at IS '作成日時';

COMMENT ON TABLE likes IS 'お酒へのいいね';
COMMENT ON COLUMN likes.id IS '一意識別子';
COMMENT ON COLUMN likes.sake_id IS 'お酒ID';
COMMENT ON COLUMN likes.token IS '一時トークン（匿名識別用, 必須）';
COMMENT ON COLUMN likes.created_at IS '作成日時';

COMMENT ON TABLE drink_logs IS '飲酒記録';
COMMENT ON COLUMN drink_logs.id IS '一意識別子';
COMMENT ON COLUMN drink_logs.user_id IS 'ユーザーID';
COMMENT ON COLUMN drink_logs.sake_id IS 'お酒ID';
COMMENT ON COLUMN drink_logs.drank_at IS '飲んだ日時';
COMMENT ON COLUMN drink_logs.created_at IS '作成日時';

COMMENT ON TABLE bookmarks IS 'ブックマーク';
COMMENT ON COLUMN bookmarks.id IS '一意識別子';
COMMENT ON COLUMN bookmarks.user_id IS 'ユーザーID';
COMMENT ON COLUMN bookmarks.sake_id IS 'お酒ID';
COMMENT ON COLUMN bookmarks.created_at IS '作成日時';

COMMENT ON TABLE match_requests IS 'マッチングリクエスト';
COMMENT ON COLUMN match_requests.id IS '一意識別子';
COMMENT ON COLUMN match_requests.from_user_id IS 'リクエスト送信者ID';
COMMENT ON COLUMN match_requests.to_user_id IS 'リクエスト受信者ID';
COMMENT ON COLUMN match_requests.status IS 'ステータス（PENDING/ACCEPTED/REJECTED）';
COMMENT ON COLUMN match_requests.created_at IS '作成日時';
COMMENT ON COLUMN match_requests.updated_at IS '更新日時';

COMMENT ON TABLE sake_images IS 'お酒画像';
COMMENT ON COLUMN sake_images.id IS '一意識別子';
COMMENT ON COLUMN sake_images.sake_id IS 'お酒ID';
COMMENT ON COLUMN sake_images.image_key IS '画像オブジェクトキー';
COMMENT ON COLUMN sake_images.sort_order IS '表示順（0始まり）';
COMMENT ON COLUMN sake_images.created_at IS '作成日時';

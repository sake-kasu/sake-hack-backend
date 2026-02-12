-- 既存テーブルのDROP(依存関係の逆順)
DROP TABLE IF EXISTS match_requests CASCADE;
DROP TABLE IF EXISTS bookmarks CASCADE;
DROP TABLE IF EXISTS drink_logs CASCADE;
DROP TABLE IF EXISTS likes CASCADE;
DROP TABLE IF EXISTS sake_drink_styles CASCADE;
DROP TABLE IF EXISTS sakes CASCADE;
DROP TABLE IF EXISTS sake_images CASCADE;
DROP TABLE IF EXISTS drink_styles CASCADE;
DROP TABLE IF EXISTS sake_kinds CASCADE;
DROP TABLE IF EXISTS breweries CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- 既存ENUMのDROP
DROP TYPE IF EXISTS match_status;
DROP TYPE IF EXISTS user_role;
DROP TYPE IF EXISTS sake_category;

-- ENUM再定義(API spec準拠)
CREATE TYPE sake_category AS ENUM (
    'JAPANESE_SAKE', 'WHISKY', 'WINE', 'BEER',
    'SHOCHU', 'AWAMORI', 'RIQUEUR', 'SPIRITS', 'OTHER'
);
COMMENT ON TYPE sake_category IS '大分類(酒のカテゴリを定義するENUM)';

CREATE TYPE user_role AS ENUM ('SHOP', 'USER');
COMMENT ON TYPE user_role IS 'ユーザーロールを定義するENUM';

CREATE TYPE match_status AS ENUM ('PENDING', 'ACCEPTED', 'REJECTED');
COMMENT ON TYPE match_status IS 'マッチングリクエストのステータスを定義するENUM';

-- ======================
-- マスタテーブル
-- ======================

CREATE TABLE breweries (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    origin_country VARCHAR(100) NOT NULL,
    origin_region VARCHAR(100),
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_breweries_name_country UNIQUE (name, origin_country)
);
COMMENT ON TABLE breweries IS '酒造マスターテーブル';
COMMENT ON COLUMN breweries.id IS '酒造の一意識別子';
COMMENT ON COLUMN breweries.name IS '酒造名';
COMMENT ON COLUMN breweries.origin_country IS '所在国';
COMMENT ON COLUMN breweries.origin_region IS '所在地域';
COMMENT ON COLUMN breweries.latitude IS '緯度';
COMMENT ON COLUMN breweries.longitude IS '経度';

CREATE TABLE sake_kinds (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE sake_kinds IS '酒の種類(小分類)マスターテーブル';
COMMENT ON COLUMN sake_kinds.id IS '酒の小分類ID';
COMMENT ON COLUMN sake_kinds.name IS '酒の小分類の名前';

CREATE TABLE drink_styles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description VARCHAR(500),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE drink_styles IS '飲み方マスターテーブル';
COMMENT ON COLUMN drink_styles.id IS '飲み方ID';
COMMENT ON COLUMN drink_styles.name IS '飲み方の名前';
COMMENT ON COLUMN drink_styles.description IS '詳細説明';

-- ======================
-- メインテーブル
-- ======================

CREATE TABLE sake_images (
    id SERIAL PRIMARY KEY,
    object_key VARCHAR(500) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE sake_images IS '酒の画像を管理するテーブル';
COMMENT ON COLUMN sake_images.id IS '画像の一意識別子';
COMMENT ON COLUMN sake_images.object_key IS '画像のオブジェクトキー';

CREATE TABLE sakes (
    id SERIAL PRIMARY KEY,
    category sake_category NOT NULL,
    kind_id INTEGER NOT NULL REFERENCES sake_kinds(id),
    brewery_id INTEGER NOT NULL REFERENCES breweries(id),
    name VARCHAR(100) NOT NULL,
    phonetic VARCHAR(200) NOT NULL,
    abv REAL NOT NULL CHECK (abv >= 0 AND abv <= 100),
    purchase_volume REAL NOT NULL CHECK (purchase_volume >= 0),
    remaining_volume REAL NOT NULL CHECK (remaining_volume >= 0),
    memo TEXT,
    price INTEGER NOT NULL CHECK (price >= 0),
    image_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_sakes_category ON sakes(category);
CREATE INDEX idx_sakes_kind_id ON sakes(kind_id);
CREATE INDEX idx_sakes_brewery_id ON sakes(brewery_id);
CREATE INDEX idx_sakes_created_at ON sakes(created_at DESC);
COMMENT ON TABLE sakes IS '酒の基本情報を管理するテーブル';
COMMENT ON COLUMN sakes.id IS '酒の一意識別子';
COMMENT ON COLUMN sakes.category IS '大分類(酒のカテゴリ)';
COMMENT ON COLUMN sakes.kind_id IS '酒の小分類ID(sake_kinds.id)';
COMMENT ON COLUMN sakes.brewery_id IS '酒造ID(breweries.id)';
COMMENT ON COLUMN sakes.name IS '酒の商品名';
COMMENT ON COLUMN sakes.phonetic IS 'ふりがな';
COMMENT ON COLUMN sakes.abv IS 'アルコール度数(%)';
COMMENT ON COLUMN sakes.purchase_volume IS '購入時容量(mL)';
COMMENT ON COLUMN sakes.remaining_volume IS '残容量(mL)';
COMMENT ON COLUMN sakes.memo IS 'メモ';
COMMENT ON COLUMN sakes.price IS '購入時価格(円)';
COMMENT ON COLUMN sakes.image_url IS '表示用の画像URL';

-- ======================
-- 中間テーブル
-- ======================

CREATE TABLE sake_drink_styles (
    sake_id INTEGER NOT NULL REFERENCES sakes(id) ON DELETE CASCADE,
    drink_style_id INTEGER NOT NULL REFERENCES drink_styles(id) ON DELETE CASCADE,
    PRIMARY KEY (sake_id, drink_style_id)
);
CREATE INDEX idx_sake_drink_styles_sake_id ON sake_drink_styles(sake_id);
CREATE INDEX idx_sake_drink_styles_drink_style_id ON sake_drink_styles(drink_style_id);
COMMENT ON TABLE sake_drink_styles IS '酒と飲み方の中間テーブル';

-- ======================
-- ユーザー関連テーブル
-- ======================

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(100),
    role user_role NOT NULL DEFAULT 'USER',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_users_email ON users(email);
COMMENT ON TABLE users IS 'ユーザー情報を管理するテーブル';
COMMENT ON COLUMN users.id IS 'ユーザーの一意識別子';
COMMENT ON COLUMN users.email IS 'メールアドレス';
COMMENT ON COLUMN users.password_hash IS 'ハッシュ化されたパスワード';
COMMENT ON COLUMN users.display_name IS '表示名';
COMMENT ON COLUMN users.role IS 'ユーザーロール(SHOP: 店舗, USER: 一般ユーザー)';

CREATE TABLE likes (
    id SERIAL PRIMARY KEY,
    sake_id INTEGER NOT NULL REFERENCES sakes(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_likes_sake_id ON likes(sake_id);
CREATE INDEX idx_likes_token ON likes(token);
COMMENT ON TABLE likes IS '酒へのいいねを管理するテーブル';

CREATE TABLE bookmarks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sake_id INTEGER NOT NULL REFERENCES sakes(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_bookmarks_user_sake UNIQUE (user_id, sake_id)
);
CREATE INDEX idx_bookmarks_user_id ON bookmarks(user_id);
CREATE INDEX idx_bookmarks_sake_id ON bookmarks(sake_id);
COMMENT ON TABLE bookmarks IS 'ユーザーのブックマークを管理するテーブル';

CREATE TABLE drink_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sake_id INTEGER NOT NULL REFERENCES sakes(id) ON DELETE CASCADE,
    drank_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_drink_logs_user_id ON drink_logs(user_id);
CREATE INDEX idx_drink_logs_sake_id ON drink_logs(sake_id);
CREATE INDEX idx_drink_logs_drank_at ON drink_logs(drank_at);
COMMENT ON TABLE drink_logs IS 'ユーザーの飲酒記録を管理するテーブル';

CREATE TABLE match_requests (
    id SERIAL PRIMARY KEY,
    from_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    to_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status match_status NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_match_requests_from_user_id ON match_requests(from_user_id);
CREATE INDEX idx_match_requests_to_user_id ON match_requests(to_user_id);
CREATE INDEX idx_match_requests_status ON match_requests(status);
COMMENT ON TABLE match_requests IS 'ユーザー間のマッチングリクエストを管理するテーブル';

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE sake_category AS ENUM ('BEER', 'JAPANESE_SAKE', 'WINE', 'WHISKEY', 'SHOCHU', 'COCKTAIL', 'OTHER');

CREATE TABLE sake (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL,
    phonetic VARCHAR(100),
    image_id UUID,
    category sake_category NOT NULL,
    description VARCHAR(100),
    alcohol_percentage NUMERIC(4, 1) CHECK (alcohol_percentage >= 0 AND alcohol_percentage <= 100),
    volume_max INTEGER CHECK (volume_max > 0),
    volume_remain INTEGER CHECK (volume_remain >= 0 AND volume_remain <= 100),
    region VARCHAR(100),
    price INTEGER CHECK (price >= 0),
    memo VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_sake_category ON sake(category);
CREATE INDEX idx_sake_name ON sake(name);

COMMENT ON TYPE sake_category IS '大分類(酒のカテゴリを定義するENUM)';

COMMENT ON TABLE sake IS '酒の基本情報を管理するテーブル';

COMMENT ON COLUMN sake.id IS '酒の一意識別子';
COMMENT ON COLUMN sake.name IS '酒の商品名';
COMMENT ON COLUMN sake.phonetic IS 'ふりがな';
COMMENT ON COLUMN sake.image_id IS '画像ID';
COMMENT ON COLUMN sake.category IS '大分類(酒のカテゴリ、(例: ビール、日本酒、ワイン))';
COMMENT ON COLUMN sake.description IS '酒の詳細情報';
COMMENT ON COLUMN sake.alcohol_percentage IS 'アルコール度数(%)';
COMMENT ON COLUMN sake.volume_max IS '最大容量(mL)';
COMMENT ON COLUMN sake.volume_remain IS '残り容量(%)';
COMMENT ON COLUMN sake.region IS '産地';
COMMENT ON COLUMN sake.price IS '価格(円)';
COMMENT ON COLUMN sake.memo IS 'メモ';
COMMENT ON COLUMN sake.created_at IS '作成日時';
COMMENT ON COLUMN sake.updated_at IS '更新日時';

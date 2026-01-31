CREATE TYPE sake_category AS ENUM ('BEER', 'JAPANESE_SAKE', 'WINE', 'WHISKEY', 'SHOCHU', 'COCKTAIL', 'OTHER');

CREATE TABLE sakes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    category sake_category NOT NULL,
    subcategory VARCHAR,
    alcohol_percentage NUMERIC(4, 1) NOT NULL CHECK (alcohol_percentage >= 0 AND alcohol_percentage <= 100),
    volume INTEGER NOT NULL CHECK (volume > 0),
    origin VARCHAR(50) NOT NULL,
    price INTEGER NOT NULL CHECK (price >= 0) DEFAULT 0,
    stock INTEGER NOT NULL CHECK (stock >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_sakes_category ON sakes(category);
CREATE INDEX idx_sakes_name ON sakes(name);

COMMENT ON TYPE sake_category IS '大分類(酒のカテゴリを定義するENUM)';

COMMENT ON TABLE sakes IS '酒の基本情報を管理するテーブル';

COMMENT ON COLUMN sakes.id IS '酒の一意識別子';
COMMENT ON COLUMN sakes.name IS '酒の商品名';
COMMENT ON COLUMN sakes.category IS '大分類(酒のカテゴリ、例)ビール、日本酒、ワイン)';
COMMENT ON COLUMN sakes.subcategory IS '小分類(酒のサブカテゴリ(例：IPA、純米大吟醸等)';
COMMENT ON COLUMN sakes.alcohol_percentage IS 'アルコール度数(%)';
COMMENT ON COLUMN sakes.volume IS '容量(mL)';
COMMENT ON COLUMN sakes.origin IS '産地';
COMMENT ON COLUMN sakes.price IS '価格(円)';
COMMENT ON COLUMN sakes.stock IS '在庫数';
COMMENT ON COLUMN sakes.created_at IS '作成日時';
COMMENT ON COLUMN sakes.updated_at IS '更新日時';

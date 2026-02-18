-- 0009で作成したテーブルをDROP
DROP TABLE IF EXISTS match_requests CASCADE;
DROP TABLE IF EXISTS drink_logs CASCADE;
DROP TABLE IF EXISTS bookmarks CASCADE;
DROP TABLE IF EXISTS likes CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS sake_drink_styles CASCADE;
DROP TABLE IF EXISTS sakes CASCADE;
DROP TABLE IF EXISTS sake_images CASCADE;
DROP TABLE IF EXISTS drink_styles CASCADE;
DROP TABLE IF EXISTS sake_kinds CASCADE;
DROP TABLE IF EXISTS breweries CASCADE;

DROP TYPE IF EXISTS match_status;
DROP TYPE IF EXISTS user_role;
DROP TYPE IF EXISTS sake_category;

-- 元のスキーマを復元(0001-0008の内容)
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE sake_category AS ENUM ('BEER', 'JAPANESE_SAKE', 'WINE', 'WHISKEY', 'SHOCHU', 'COCKTAIL', 'OTHER');
CREATE TYPE user_role AS ENUM ('SHOP', 'USER');
CREATE TYPE match_status AS ENUM ('PENDING', 'ACCEPTED', 'REJECTED');

-- 0001: sakes
CREATE TABLE sakes (
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
CREATE INDEX idx_sake_category ON sakes(category);
CREATE INDEX idx_sake_name ON sakes(name);

-- 0002: sake_images
CREATE TABLE sake_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    image_url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE INDEX idx_sake_images_created_at ON sake_images(created_at);

-- 0003: sake_image FK
ALTER TABLE sakes ADD CONSTRAINT fk_sakes_image_id FOREIGN KEY (image_id) REFERENCES sake_images(id) ON DELETE SET NULL;

-- 0004: users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(100),
    role user_role NOT NULL DEFAULT 'USER',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE INDEX idx_users_email ON users(email);

-- 0005: likes
CREATE TABLE likes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sake_id UUID NOT NULL REFERENCES sakes(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE INDEX idx_likes_sake_id ON likes(sake_id);
CREATE INDEX idx_likes_token ON likes(token);

-- 0006: drink_logs
CREATE TABLE drink_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sake_id UUID NOT NULL REFERENCES sakes(id) ON DELETE CASCADE,
    drank_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE INDEX idx_drink_logs_user_id ON drink_logs(user_id);
CREATE INDEX idx_drink_logs_sake_id ON drink_logs(sake_id);
CREATE INDEX idx_drink_logs_drank_at ON drink_logs(drank_at);

-- 0007: bookmarks
CREATE TABLE bookmarks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sake_id UUID NOT NULL REFERENCES sakes(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT uq_bookmarks_user_sake UNIQUE (user_id, sake_id)
);
CREATE INDEX idx_bookmarks_user_id ON bookmarks(user_id);
CREATE INDEX idx_bookmarks_sake_id ON bookmarks(sake_id);

-- 0008: match_requests
CREATE TABLE match_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    to_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status match_status NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE INDEX idx_match_requests_from_user_id ON match_requests(from_user_id);
CREATE INDEX idx_match_requests_to_user_id ON match_requests(to_user_id);
CREATE INDEX idx_match_requests_status ON match_requests(status);

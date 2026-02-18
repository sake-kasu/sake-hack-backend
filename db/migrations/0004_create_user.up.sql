CREATE TYPE user_role AS ENUM ('SHOP', 'USER');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(100),
    role user_role NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_users_email ON users(email);

COMMENT ON TYPE user_role IS 'ユーザーロールを定義するENUM';

COMMENT ON TABLE users IS 'ユーザー情報を管理するテーブル';

COMMENT ON COLUMN users.id IS 'ユーザーの一意識別子';
COMMENT ON COLUMN users.email IS 'メールアドレス';
COMMENT ON COLUMN users.password_hash IS 'ハッシュ化されたパスワード';
COMMENT ON COLUMN users.display_name IS '表示名';
COMMENT ON COLUMN users.role IS 'ユーザーロール(SHOP: 店舗, USER: 一般ユーザー)';
COMMENT ON COLUMN users.created_at IS '作成日時';
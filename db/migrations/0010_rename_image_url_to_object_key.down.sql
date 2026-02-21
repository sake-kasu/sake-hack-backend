ALTER TABLE sakes RENAME COLUMN object_key TO image_url;
COMMENT ON COLUMN sakes.image_url IS '表示用の画像URL';
CREATE TABLE sake_images (
    id SERIAL PRIMARY KEY,
    object_key VARCHAR(500) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

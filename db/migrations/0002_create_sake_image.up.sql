CREATE TABLE sake_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    object_key VARCHAR(500) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_sake_images_created_at ON sake_images(created_at);

COMMENT ON TABLE sake_images IS '酒の画像を管理するテーブル';

COMMENT ON COLUMN sake_images.id IS '画像の一意識別子';
COMMENT ON COLUMN sake_images.image_url IS 'S3に配置される画像のパス';
COMMENT ON COLUMN sake_images.created_at IS '作成日時';
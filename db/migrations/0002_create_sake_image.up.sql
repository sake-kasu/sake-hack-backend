CREATE TABLE sake_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    image_url VARCHAR(500) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_sake_images_created_at ON sake_images(created_at);

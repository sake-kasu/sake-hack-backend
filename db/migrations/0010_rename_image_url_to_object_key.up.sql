DROP TABLE IF EXISTS sake_images CASCADE;
ALTER TABLE sakes RENAME COLUMN image_url TO object_key;
COMMENT ON COLUMN sakes.object_key IS 'S3/RustFSのオブジェクトキー';

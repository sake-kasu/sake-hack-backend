ALTER TABLE sakes
ADD CONSTRAINT fk_sake_image_id
FOREIGN KEY (image_id) REFERENCES sake_images(id) ON DELETE SET NULL;

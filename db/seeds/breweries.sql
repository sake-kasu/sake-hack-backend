-- 酒蔵マスターデータ
INSERT INTO breweries (name, origin_country, origin_region, latitude, longitude) VALUES
('旭酒造', '日本', '山口県', 34.1234, 131.4567),
('久保田酒造', '日本', '新潟県', 37.9023, 139.0234),
('獺祭酒造', '日本', '山口県', 34.1500, 131.5000)
ON CONFLICT (name, origin_country) DO NOTHING;

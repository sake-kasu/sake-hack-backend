-- 酒の種類（小分類）マスターデータ
INSERT INTO sake_kinds (name) VALUES
('純米大吟醸'),
('純米吟醸'),
('本醸造')
ON CONFLICT (name) DO NOTHING;

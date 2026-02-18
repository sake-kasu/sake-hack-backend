-- 飲み方マスターデータ
INSERT INTO drink_styles (name, description) VALUES
('冷酒', '冷やして飲む'),
('常温', '常温で飲む'),
('熱燗', '温めて飲む')
ON CONFLICT (name) DO NOTHING;

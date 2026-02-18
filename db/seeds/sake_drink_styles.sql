-- 酒と飲み方の関連付け
-- 依存: sakes, drink_styles
INSERT INTO sake_drink_styles (sake_id, drink_style_id) VALUES
(1, 1), -- 獺祭は冷酒
(1, 2), -- 獺祭は常温
(2, 1), -- 久保田は冷酒
(2, 2), -- 久保田は常温
(2, 3), -- 久保田は熱燗
(3, 1) -- 獺祭45は冷酒
ON CONFLICT (sake_id, drink_style_id) DO NOTHING;

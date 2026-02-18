-- 酒データ
-- 依存: breweries, sake_kinds
INSERT INTO sakes (category, kind_id, brewery_id, name, phonetic, abv, purchase_volume, remaining_volume, memo, price, image_url)
SELECT 'JAPANESE_SAKE', 1, 1, '獺祭 純米大吟醸 磨き二割三分', 'だっさい じゅんまいだいぎんじょう みがきにわりさんぶ', 16.0, 720, 720, '最高級の日本酒', 5000, 'https://example.com/dassai.jpg'
WHERE NOT EXISTS (SELECT 1 FROM sakes WHERE name = '獺祭 純米大吟醸 磨き二割三分' AND brewery_id = 1);

INSERT INTO sakes (category, kind_id, brewery_id, name, phonetic, abv, purchase_volume, remaining_volume, memo, price, image_url)
SELECT 'JAPANESE_SAKE', 2, 2, '久保田 萬寿', 'くぼた まんじゅ', 15.5, 1800, 1800, '新潟の名酒', 8000, 'https://example.com/kubota.jpg'
WHERE NOT EXISTS (SELECT 1 FROM sakes WHERE name = '久保田 萬寿' AND brewery_id = 2);

INSERT INTO sakes (category, kind_id, brewery_id, name, phonetic, abv, purchase_volume, remaining_volume, memo, price, image_url)
SELECT 'JAPANESE_SAKE', 1, 3, '獺祭 純米大吟醸 45', 'だっさい じゅんまいだいぎんじょう よんじゅうご', 16.0, 720, 600, 'スタンダードな獺祭', 2500, NULL
WHERE NOT EXISTS (SELECT 1 FROM sakes WHERE name = '獺祭 純米大吟醸 45' AND brewery_id = 3);

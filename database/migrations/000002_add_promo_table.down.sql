-- Hapus kolom dari transactions dulu (karena foreign key)
ALTER TABLE transactions 
DROP COLUMN IF EXISTS final_amount,
DROP COLUMN IF EXISTS discount_amount,
DROP COLUMN IF EXISTS promo_id;

-- Hapus tabel promos
DROP TABLE IF EXISTS promos;
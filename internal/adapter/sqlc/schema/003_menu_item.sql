-- +goose Up
CREATE TABLE menu_items (
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    name VARCHAR(100) NOT NULL,
    description TEXT DEFAULT '',
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_menu_items_category_id ON menu_items(category_id);
CREATE INDEX idx_menu_items_name ON menu_items(name);
CREATE INDEX idx_menu_items_is_active ON menu_items(is_active);
CREATE INDEX idx_menu_items_category_active ON menu_items(category_id, is_active);
CREATE INDEX idx_menu_items_name_search ON menu_items USING gin(to_tsvector('english', name || ' ' || description));

-- Insert sample menu items
INSERT INTO menu_items (category_id, name, description, price) VALUES
-- ของคาว (category_id = 1)
(1, 'ข้าวผัดกุ้ง', 'ข้าวผัดกุ้งสดใส่ไข่ หอมหวาน', 120.00),
(1, 'ผัดไทยกุ้งสด', 'ผัดไทยแท้รสชาติดั้งเดิม', 150.00),
(1, 'ต้มยำกุ้งน้ำข้น', 'ต้มยำกุ้งรสจัดจ้าน เผ็ดมาก', 180.00),
(1, 'แกงเขียวหวานไก่', 'แกงเขียวหวานไก่ใส่มะเขือ', 160.00),
(1, 'ยำวุ้นเส้น', 'ยำวุ้นเส้นทะเลรสเปรี้ยว', 130.00),

-- ของหวาน (category_id = 2)
(2, 'ข้าวเหนียวมะม่วง', 'ข้าวเหนียวมะม่วงน้ำดอกไม้', 80.00),
(2, 'ทับทิมกรอบ', 'ทับทิมกรอบกะทิสด', 60.00),
(2, 'บัวลอยไข่หวาน', 'บัวลอยไข่หวานน้ำกะทิ', 50.00),
(2, 'กล้วยบวชชี', 'กล้วยบวชชีกะทิข้นหวาน', 45.00),
(2, 'ฟักทองแกงบวด', 'ฟักทองแกงบวดกะทิ', 55.00),

-- ของทานเล่น (category_id = 3)
(3, 'ปอเปี๊ยะทอด', 'ปอเปี๊ยะทอดกรอบไส้ผัก', 40.00),
(3, 'เต้าหู้ทอด', 'เต้าหู้ทอดกรอบซอสหวาน', 35.00),
(3, 'ไส้กรอกอีสาน', 'ไส้กรอกอีสานรสเปรี้ยว', 60.00),
(3, 'หมูปิ้ง', 'หมูปิ้งไฟฟ้าซอสแจ่ว', 80.00),
(3, 'ลาบหมู', 'ลาบหมูรสจัด เผ็ดมาก', 90.00),

-- โรตี (category_id = 4)
(4, 'โรตีกล้วยไข่', 'โรตีกล้วยไข่หวานมัน', 45.00),
(4, 'โรตีนมสด', 'โรตีนมสดนมข้นหวาน', 40.00),
(4, 'โรตีช็อกโกแลต', 'โรตีช็อกโกแลตกรอบหวาน', 50.00),
(4, 'โรตีสังขยา', 'โรตีสังขยาใบเตยหอม', 55.00),
(4, 'โรตีมะตะบะ', 'โรตีมะตะบะแกงไก่', 70.00);

-- +goose Down
DROP INDEX IF EXISTS idx_menu_items_name_search;
DROP INDEX IF EXISTS idx_menu_items_category_active;
DROP INDEX IF EXISTS idx_menu_items_is_active;
DROP INDEX IF EXISTS idx_menu_items_name;
DROP INDEX IF EXISTS idx_menu_items_category_id;
DROP TABLE IF EXISTS menu_items;
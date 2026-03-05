# CorePOS

Core Point of Sale System.


-- ==========================================
-- 1. สร้างประเภทข้อมูล (ENUM Types)
-- ==========================================
CREATE TYPE order_status AS ENUM ('pending', 'completed', 'void', 'refund');
CREATE TYPE payment_status_enum AS ENUM ('unpaid', 'partial', 'paid');

-- ==========================================
-- 2. สร้างฟังก์ชัน Trigger สำหรับอัปเดตเวลาอัตโนมัติ
-- ==========================================
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ==========================================
-- 3. สร้างตารางทั้งหมด (Tables Creation)
-- ==========================================

-- 3.1 ตารางร้านค้า (ผู้เช่าระบบ SaaS)
CREATE TABLE stores (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    plan_type VARCHAR(50) DEFAULT 'free',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3.2 ตารางประวัติการต่ออายุ SaaS
CREATE TABLE subscription_history (
    id SERIAL PRIMARY KEY,
    store_id INTEGER NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    plan_name VARCHAR(50) NOT NULL,
    amount_paid DECIMAL(10, 2) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    payment_status VARCHAR(50) DEFAULT 'success',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3.3 ตารางผู้ใช้งาน/พนักงาน
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    store_id INTEGER NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    username VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'cashier',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3.4 ตารางหมวดหมู่สินค้า
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    store_id INTEGER NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3.5 ตารางสินค้า
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    store_id INTEGER NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    barcode VARCHAR(100),
    cost_price DECIMAL(10, 2) DEFAULT 0.00,
    price DECIMAL(10, 2) NOT NULL,
    stock_quantity INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL -- Soft Delete
);

-- 3.6 ตารางประวัติคลังสินค้า (เข้า-ออก)
CREATE TABLE inventory_movements (
    id SERIAL PRIMARY KEY,
    store_id INTEGER NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    movement_type VARCHAR(50) NOT NULL,
    quantity_changed INTEGER NOT NULL,
    reference_id INTEGER,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3.7 ตารางใบสั่งซื้อ/บิลการขาย
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    store_id INTEGER NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id),
    total_amount DECIMAL(10, 2) NOT NULL,
    status order_status DEFAULT 'completed',
    payment_status payment_status_enum DEFAULT 'unpaid',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3.8 ตารางรายละเอียดในบิล (สินค้าที่ขาย)
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(id) ON DELETE SET NULL,
    quantity INTEGER NOT NULL,
    unit_price DECIMAL(10, 2) NOT NULL,
    cost_price_snapshot DECIMAL(10, 2) NOT NULL,
    subtotal DECIMAL(10, 2) NOT NULL
);

-- 3.9 ตารางการชำระเงิน (รองรับการแบ่งจ่าย)
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    amount DECIMAL(10, 2) NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    paid_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ==========================================
-- 4. เปิดใช้งาน Triggers
-- ==========================================
CREATE TRIGGER set_timestamp_products
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp_orders
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- ==========================================
-- 5. สร้างดัชนี (Indexes) เพื่อเพิ่มความเร็ว
-- ==========================================
CREATE INDEX idx_products_store_id ON products(store_id);
CREATE INDEX idx_products_deleted_at ON products(deleted_at);
CREATE INDEX idx_orders_store_id ON orders(store_id);
CREATE INDEX idx_orders_created_at ON orders(created_at);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_inventory_product_id ON inventory_movements(product_id);
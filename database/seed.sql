-- Battery ERP System Seed Data
-- Initialize the system with default data

USE battery_erp;

-- Insert default admin user
-- Password: admin123 (hashed with bcrypt)
INSERT INTO users (username, password, real_name, role, is_active) VALUES 
('admin', '$2a$10$9VgzQX1KXzA8yLX7ZGrD.egJh7YfGg8z7YYJy8x0c7Hf8Vv6FjQ.K', '系统管理员', 'super_admin', TRUE),
('operator', '$2a$10$9VgzQX1KXzA8yLX7ZGrD.egJh7YfGg8z7YYJy8x0c7Hf8Vv6FjQ.K', '操作员', 'normal', TRUE)
ON DUPLICATE KEY UPDATE username=username;

-- Insert default battery categories
INSERT INTO battery_categories (name, description, unit_price, is_active) VALUES 
('铅酸电池', '废旧铅酸蓄电池', 8.50, TRUE),
('锂电池', '废旧锂离子电池', 15.00, TRUE),
('镍氢电池', '废旧镍氢电池', 12.00, TRUE),
('镍镉电池', '废旧镍镉电池', 10.00, TRUE),
('碱性电池', '废旧碱性电池', 6.00, TRUE)
ON DUPLICATE KEY UPDATE name=name;

-- Initialize inventory records for all categories
INSERT INTO inventory (category_id, current_weight_kg, created_at, updated_at)
SELECT id, 0, NOW(), NOW() FROM battery_categories
ON DUPLICATE KEY UPDATE category_id=category_id;

-- Insert sample inbound order (optional - for testing)
INSERT INTO inbound_orders (order_no, supplier_name, total_amount, status, notes, created_by) VALUES 
('IN-20250629-0001', '测试供应商', 850.00, 'completed', '系统初始化测试数据', 1)
ON DUPLICATE KEY UPDATE order_no=order_no;

-- Insert sample inbound order items (optional - for testing)
INSERT INTO inbound_order_items (order_id, category_id, gross_weight, tare_weight, net_weight, unit_price, sub_total) 
SELECT 
    (SELECT id FROM inbound_orders WHERE order_no = 'IN-20250629-0001'),
    (SELECT id FROM battery_categories WHERE name = '铅酸电池'),
    105.000,
    5.000,
    100.000,
    8.50,
    850.00
WHERE EXISTS (SELECT 1 FROM inbound_orders WHERE order_no = 'IN-20250629-0001')
AND EXISTS (SELECT 1 FROM battery_categories WHERE name = '铅酸电池')
ON DUPLICATE KEY UPDATE order_id=order_id;

-- Update inventory for the sample data
UPDATE inventory 
SET current_weight_kg = 100.000, last_inbound_at = NOW(), updated_at = NOW()
WHERE category_id = (SELECT id FROM battery_categories WHERE name = '铅酸电池');

-- Show summary of initialized data
SELECT 'Users' as table_name, COUNT(*) as count FROM users
UNION ALL
SELECT 'Battery Categories', COUNT(*) FROM battery_categories
UNION ALL  
SELECT 'Inventory Records', COUNT(*) FROM inventory
UNION ALL
SELECT 'Sample Orders', COUNT(*) FROM inbound_orders;
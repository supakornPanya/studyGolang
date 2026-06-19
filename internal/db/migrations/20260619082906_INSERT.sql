-- +goose Up
-- Insert Tags
INSERT INTO tags (id, name) VALUES
(1, 'Electronics'),
(2, 'Gadgets'),
(3, 'Home Office'),
(4, 'Kitchen'),
(5, 'Smart Home');

-- Insert Categories
INSERT INTO categories (id, name, description, parent_id) VALUES
(1, 'Electronics', 'Electronic devices and accessories', 0),
(2, 'Computers & Laptops', 'Personal computers, laptops, and tablets', 1),
(3, 'Home Appliances', 'Appliances for home and kitchen', 0);

-- Insert Products
-- Note: 'created_by' is set to 1, representing a seller's user ID.
-- 'category' and 'tag' are JSON arrays of IDs.
INSERT INTO products (id, sku, name, description, price, stock, category, tag, created_by) VALUES
(1, 'PROD-001', 'Wireless Noise-Canceling Headphones', 'High-quality wireless headphones with active noise canceling and 30-hour battery life.', 199.99, 50, '[1, 2]', '[1, 2]', 1),
(2, 'PROD-002', 'Ergonomic Mechanical Keyboard', 'RGB backlit mechanical keyboard with tactile switches and wrist rest.', 89.99, 120, '[1, 2]', '[3]', 1),
(3, 'PROD-003', 'Smart Coffee Maker', 'Wi-Fi enabled coffee maker with automatic brewing and temperature control.', 129.99, 30, '[3]', '[4, 5]', 1);

-- Insert Cart Items (matching internal/domain/entity/cart_items.go)
INSERT INTO cart_items (user_id, list_product_id, created_at, updated_at) VALUES
(1, '[1, 2]', NOW(), NOW()),
(2, '[3]', NOW(), NOW());

-- Reset serial sequences so future AUTO_INCREMENT inserts work correctly
SELECT setval(pg_get_serial_sequence('tags', 'id'), COALESCE(MAX(id), 1)) FROM tags;
SELECT setval(pg_get_serial_sequence('categories', 'id'), COALESCE(MAX(id), 1)) FROM categories;
SELECT setval(pg_get_serial_sequence('products', 'id'), COALESCE(MAX(id), 1)) FROM products;

-- +goose Down
DELETE FROM cart_items WHERE user_id IN (1, 2);
DELETE FROM products WHERE id IN (1, 2, 3);
DELETE FROM categories WHERE id IN (1, 2, 3);
DELETE FROM tags WHERE id IN (1, 2, 3, 4, 5);

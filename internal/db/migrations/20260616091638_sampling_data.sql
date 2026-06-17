-- +goose Up
-- Insert Items
INSERT INTO
    items (sku, name, quantity, price)
VALUES (
        'SKU001',
        'Product A',
        10,
        99.99
    );

INSERT INTO
    items (sku, name, quantity, price)
VALUES (
        'SKU002',
        'Product B',
        20,
        199.99
    );

INSERT INTO
    items (sku, name, quantity, price)
VALUES (
        'SKU003',
        'Product C',
        30,
        299.99
    );

INSERT INTO
    items (sku, name, quantity, price)
VALUES (
        'SKU004',
        'Product D',
        40,
        399.99
    );

-- +goose Down
DELETE FROM items
WHERE
    sku IN (
        'SKU001',
        'SKU002',
        'SKU003',
        'SKU004'
    );
-- +migrate Up
INSERT INTO shopping_categories (name, created_at, updated_at)
VALUES 
    ('foods', NOW(), NOW()),
    ('daily', NOW(), NOW()),
    ('hygiene', NOW(), NOW()),
    ('pet', NOW(), NOW()),
    ('healthcare', NOW(), NOW()),
    ('miscellaneous', NOW(), NOW()),
    ('hobby', NOW(), NOW());

-- +migrate Down
DELETE FROM shopping_categories 
WHERE name IN ('foods', 'daily', 'hygiene', 'pet', 'healthcare', 'miscellaneous', 'hobby'); 

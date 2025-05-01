-- +migrate Up
CREATE TABLE IF NOT EXISTS personal_shopping_items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(10) NOT NULL,
    category_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (category_id) REFERENCES shopping_categories(id)
);

-- +migrate Down
DROP TABLE IF EXISTS personal_shopping_items;

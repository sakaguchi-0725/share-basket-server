-- +migrate Up
CREATE TABLE IF NOT EXISTS personal_shopping_items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(100) NOT NULL,
    category_id INTEGER NOT NULL,
    account_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_personal_items_categories FOREIGN KEY (category_id) REFERENCES categories(id),
    CONSTRAINT fk_personal_items_accounts FOREIGN KEY (account_id) REFERENCES accounts(id)
);

-- +migrate Down
DROP TABLE IF EXISTS personal_shopping_items;

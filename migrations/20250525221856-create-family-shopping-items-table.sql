-- +migrate Up
CREATE TABLE IF NOT EXISTS family_shopping_items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(100) NOT NULL,
    category_id INTEGER NOT NULL,
    family_id UUID NOT NULL,
    created_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_family_items_categories FOREIGN KEY (category_id) REFERENCES categories(id),
    CONSTRAINT fk_family_items_families FOREIGN KEY (family_id) REFERENCES families(id),
    CONSTRAINT fk_family_items_accounts FOREIGN KEY (created_by) REFERENCES accounts(id)
);

-- +migrate Down
DROP TABLE IS EXISTS family_shopping_items;

-- +migrate Up
CREATE TABLE IF NOT EXISTS family_members (
    id SERIAL PRIMARY KEY,
    family_id UUID NOT NULL,
    account_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_family_members_families FOREIGN KEY (family_id) REFERENCES families(id),
    CONSTRAINT fk_family_members_accounts FOREIGN KEY (account_id) REFERENCES accounts(id)
);

-- +migrate Down
DROP TABLE NOT EXISTS family_members;

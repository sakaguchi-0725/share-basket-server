-- +migrate Up
INSERT INTO categories (name) VALUES ('foods');          -- 食料品
INSERT INTO categories (name) VALUES ('daily');          -- 日用品
INSERT INTO categories (name) VALUES ('clothing');       -- 衣料品
INSERT INTO categories (name) VALUES ('cosmetics');      -- 化粧品
INSERT INTO categories (name) VALUES ('pet');            -- ペット用品
INSERT INTO categories (name) VALUES ('baby');           -- ベビー用品
INSERT INTO categories (name) VALUES ('homeAppliances'); -- 家電
INSERT INTO categories (name) VALUES ('hobby');          -- 趣味・趣向品

-- +migrate Down
TRUNCATE TABLE categories RESTART IDENTITY;


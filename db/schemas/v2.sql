ALTER TABLE records
ALTER COLUMN name SET DATA TYPE VARCHAR(255) COLLATE "unicode";

ALTER TABLE items
ALTER COLUMN name TYPE VARCHAR(255);

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL COLLATE "unicode",
    color VARCHAR(30) NOT NULL,
    clerk_user_id TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp,
    updated_at TIMESTAMP DEFAULT NULL,
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TRIGGER update_updated_at BEFORE UPDATE ON categories FOR EACH ROW EXECUTE PROCEDURE update_updated_at();

CREATE INDEX categories_clerk_user_id_idx ON categories (clerk_user_id);

ALTER TABLE items
ADD COLUMN category_id INTEGER;

CREATE INDEX items_category_id_idx ON items (category_id);
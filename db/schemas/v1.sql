CREATE OR REPLACE FUNCTION update_updated_at() RETURNS TRIGGER AS
$$
    BEGIN NEW.updated_at = now();
    RETURN NEW;
    END;
$$ language 'plpgsql';

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    clerk_user_id TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp,
    updated_at TIMESTAMP DEFAULT NULL,
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TRIGGER update_updated_at BEFORE UPDATE ON items FOR EACH ROW EXECUTE PROCEDURE update_updated_at();

CREATE INDEX items_clerk_user_id_idx ON items (clerk_user_id);

CREATE TABLE records (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    value TEXT NOT NULL,
    item_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp,
    updated_at TIMESTAMP DEFAULT NULL,
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TRIGGER update_updated_at BEFORE UPDATE ON records FOR EACH ROW EXECUTE PROCEDURE update_updated_at();

CREATE INDEX records_item_id_idx ON records (item_id);

ALTER TABLE items
ALTER COLUMN name
SET DATA TYPE text COLLATE "unicode";
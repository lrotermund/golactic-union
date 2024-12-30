-- +migrate Up
CREATE table space_ships (
    id BLOB PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    intact INTEGER DEFAULT 1 NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE space_ships;

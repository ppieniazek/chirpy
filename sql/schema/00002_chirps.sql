-- +goose Up
CREATE TABLE chirps (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    body TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS chirps;

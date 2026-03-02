-- +goose Up
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    feed_id UUID NOT NULL REFERENCES feeds(id) on delete cascade,
    user_id UUID NOT NULL REFERENCES users(id) on delete cascade,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT feedid_userid_unique
        UNIQUE (feed_id, user_id)
);

-- +goose Down
DROP TABLE feed_follows;
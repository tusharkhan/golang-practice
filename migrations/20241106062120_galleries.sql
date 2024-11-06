-- +goose Up
-- +goose StatementBegin
CREATE TABLE galleries(
    id SERIAL PRIMARY KEY,
    user_id int REFERENCES users (id),
    title TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) ;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE galleries;
-- +goose StatementEnd

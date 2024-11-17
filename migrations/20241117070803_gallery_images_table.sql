-- +goose Up
-- +goose StatementBegin
CREATE TABLE gallery_images(
    id SERIAL PRIMARY KEY,
    gallery_id int REFERENCES galleries (id),
    real_name TEXT,
    generate_name TEXT,
    file_size INT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) ;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE gallery_images;
-- +goose StatementEnd

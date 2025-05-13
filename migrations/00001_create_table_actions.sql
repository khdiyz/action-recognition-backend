-- +goose Up
-- +goose StatementBegin
CREATE TABLE actions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    video_id varchar(64) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE actions;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE actions (
    id SERIAL PRIMARY KEY,
    video_url VARCHAR(255) NOT NULL,
    predicted_actions varchar[],
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE actions;
-- +goose StatementEnd

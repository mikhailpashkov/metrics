-- +goose Up
CREATE TABLE metrics
(
    id    BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    ts    TIMESTAMP NOT NULL,
    type  TEXT      NOT NULL,
    name  TEXT      NOT NULL,
    delta BIGINT,
    value DOUBLE PRECISION
);

-- +goose Down
DROP TABLE metrics

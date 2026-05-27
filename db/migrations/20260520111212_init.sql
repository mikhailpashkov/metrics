-- +goose Up
CREATE TABLE metrics
(
    id    BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    ts    TIMESTAMP    NOT NULL,
    type  VARCHAR(64)  NOT NULL,
    name  VARCHAR(256) NOT NULL,
    delta BIGINT,
    value DOUBLE PRECISION
);

CREATE INDEX idx_metrics_name ON metrics (name);

-- +goose Down
DROP TABLE metrics
DROP INDEX IF EXISTS idx_metrics_name;

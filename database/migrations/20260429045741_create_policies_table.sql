-- +goose Up
-- +goose StatementBegin
CREATE TABLE policies (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    ptype varchar(2) NOT NULL CHECK (ptype IN ('p', 'g', 'g2', 'g3')),
    v0 varchar(100) NULL,
    v1 varchar(100) NULL,
    v2 varchar(100) NULL,
    v3 varchar(100) NULL,
    v4 varchar(100) NULL,
    v5 varchar(100) NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_policies_ptype ON policies(ptype);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS policies;
-- +goose StatementEnd
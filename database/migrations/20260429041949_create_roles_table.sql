-- +goose Up
-- +goose StatementBegin
CREATE TABLE roles (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    slug varchar(64) NOT NULL UNIQUE,
    lv smallint NOT NULL DEFAULT 0 CHECK (lv >= 0 AND lv <= 100),
    name varchar(64) NOT NULL,
    permissions jsonb NOT NULL DEFAULT '{"guest": true}'::jsonb,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
    -- CONSTRAINT roles_slug_unique UNIQUE (slug)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd
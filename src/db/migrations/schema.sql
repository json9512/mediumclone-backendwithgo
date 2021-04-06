-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email varchar(255) UNIQUE,
    pwd varchar(255) ,
    token_expires_in bigint,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY UNIQUE,
    author varchar(255),
    document text,
    comments text,
    likes int,
    tags text[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS tags_index ON posts USING GIN("tags");

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION update_timestamp() RETURNS TRIGGER AS $$
    BEGIN
        NEW.updated_at = NOW();
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_timestamp() RETURNS TRIGGER AS $$
    BEGIN
        NEW.deleted_at = NOW();
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER update_stamp BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE PROCEDURE update_timestamp();

CREATE TRIGGER update_stamp BEFORE UPDATE ON posts
    FOR EACH ROW EXECUTE PROCEDURE update_timestamp();

CREATE TRIGGER delete_stamp BEFORE UPDATE ON posts
    FOR EACH ROW EXECUTE PROCEDURE delete_timestamp();
-- +migrate StatementEnd

-- +migrate Down
DROP TABLE users;
DROP TABLE posts;
DROP FUNCTION update_timestamp();
DROP FUNCTION delete_timestamp();

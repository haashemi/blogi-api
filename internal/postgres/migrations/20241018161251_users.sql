-- +goose Up
-- +goose StatementBegin
CREATE TABLE users(
    id          bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    full_name   varchar(256)    NOT NULL,
    username    varchar(32)     NOT NULL    CHECK (username ~ '^(?=.{4,32}$)(?![_.])(?!.*[_.]{2})[a-zA-Z0-9._]+(?<![_.])$'),
    about_me    varchar(2048)           ,
    password    text            NOT NULL,
    is_admin    boolean         NOT NULL    DEFAULT false,
    is_banned   boolean         NOT NULL    DEFAULT false,
    created_at  timestamptz     NOT NULL    DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

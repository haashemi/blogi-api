-- +goose Up
-- +goose StatementBegin
CREATE TABLE blogs(
    id          bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY ,
    author_id   bigint          NOT NULL    REFERENCES users(id),
    title       varchar(1024)   NOT NULL                        ,
    summary     varchar(2048)   NOT NULL                        ,
    content     text            NOT NULL                        ,
    created_at  timestamptz     NOT NULL    DEFAULT NOW()       ,
    updated_at  timestamptz     NOT NULL    DEFAULT NOW()       ,
    removed_at  timestamptz
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS blogs;
-- +goose StatementEnd

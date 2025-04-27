-- +goose Up
-- +goose StatementBegin
create table users (
  id UUID primary key,
  created_at timestamp not null,
  updated_at timestamp not null,
  name varchar(255) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users;
-- +goose StatementEnd


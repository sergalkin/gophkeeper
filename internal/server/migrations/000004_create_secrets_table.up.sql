create table secrets
(
    id         bigserial primary key,
    user_id    uuid      not null,
    type_id    bigserial not null,
    title      text      not null,
    content    bytea     not null,
    created_at TIMESTAMPTZ default now(),
    updated_at TIMESTAMPTZ default now(),
    deleted_at TIMESTAMPTZ default null,

    constraint fk_type_id foreign key (type_id) references secret_types (id) on delete cascade,
    constraint fk_user_id foreign key (user_id) references users (id) on delete cascade
);
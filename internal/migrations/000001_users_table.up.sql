create TABLE users
(
    id         UUID        DEFAULT gen_random_uuid() not null unique,
    email      TEXT                                  not null,
    password   TEXT                                  not null,
    created_at TIMESTAMPTZ default now()
);
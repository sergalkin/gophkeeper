create table secret_types
(
    id    bigserial primary key,
    title text not null
);

insert into secret_types (title)
values ('login/pass'),
       ('text'),
       ('binary'),
       ('card')
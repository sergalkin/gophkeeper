create table types
(
    id    bigserial primary key,
    title text not null
);

insert into types (title)
values ('login/pass'),
       ('text'),
       ('binary'),
       ('card')
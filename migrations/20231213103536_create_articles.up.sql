create table if not exists articles
(
    id           serial primary key,
    source_id    bigint unsigned not null,
    title        varchar(255)    not null,
    link         varchar(255)    not null unique,
    summary      TEXT            not null,
    published_at timestamp       not null,
    created_at   timestamp       not null default now(),
    posted_at    timestamp       null
);
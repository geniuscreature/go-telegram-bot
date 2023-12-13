create table articles (
    id serial primary key,
    source_id int not null,
    title varchar(255) not null,
    link varchar(255) not null,
    summary TEXT not null,
    published_at timestamp not null,
    created_at timestamp not null,
    updated_at timestamp
);
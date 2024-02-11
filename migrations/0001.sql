
create table shoppinglists (
    id bigserial primary key,
    uuid text not null,
    unique(uuid)
);

create table items (
    id bigserial primary key,
    uuid integer not null,
    shoppinglist_id bigint references shoppinglists on delete cascade,
    name text not null,
    completed boolean not null,
    unique(uuid)
);

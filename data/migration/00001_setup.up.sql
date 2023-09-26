Create table if not exists product (
    id bigint primary key,
    name varchar not null,
    price integer not null check(price > 0),
    created_at timestamptz not null,
    updated_at timestamptz  not null,
    deleted_at timestamptz 
)
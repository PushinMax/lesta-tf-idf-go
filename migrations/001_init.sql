CREATE TABLE IF NOT EXISTS users (
    id uuid primary key,
    login varchar(20),
    password_hash text,
    token_hash text,
    created_at timestamptz,
);








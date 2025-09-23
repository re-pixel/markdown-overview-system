create table if not exists user_sessions (
    id serial primary key,
    user_id int not null references users(id) on delete cascade,
    session_token varchar(255) unique not null,
    created_at timestamp default current_timestamp,
    expires_at timestamp not null
);
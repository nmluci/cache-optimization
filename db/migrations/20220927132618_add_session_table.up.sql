create table user_sessions (
    session_id varchar(255) not null unique primary key,
    user_id bigint(20) not null unique,
    expired_at bigint(20) not null
);

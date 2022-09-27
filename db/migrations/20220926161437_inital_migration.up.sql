create table userdata (
    id bigint(20) not null auto_increment primary key,
    email varchar(255) not null,
    password varchar(255) not null,
    fullname varchar(255) not null,
    created_at datetime not null default now(),
    updated_at datetime not null default now(),
    deleted_at datetime
);

create table products (
    id bigint(20) not null auto_increment primary key,
    name varchar(255) not null,
    category varchar(255) not null,
    description varchar(255) not null,
    unit_price int(20) not null,
    qty int(20) not null default 0,
    created_at datetime not null default now(),
    updated_at datetime not null default now(),
    deleted_at datetime
);

create table orders (
    id bigint(20) not null auto_increment primary key,
    userdata_id bigint(20) not null,
    order_date datetime not null default now(),
    paid_at datetime,
    created_at datetime not null default now(),
    updated_at datetime not null default now(),
    deleted_at datetime
);

create table order_details (
    id bigint(20) not null auto_increment primary key,
    order_id bigint(20) not null,
    product_id bigint(20) not null,
    qty int(20) not null,
    price int(20) not null,
    created_at datetime not null default now(),
    updated_at datetime not null default now(),
    deleted_at datetime
);
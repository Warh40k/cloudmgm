CREATE TABLE users
(
    id uuid primary key,
    name varchar(255) not null,
    username varchar(255) not null unique,
    password_hash varchar(255) not null,
    balance double precision not null default 0,
    created timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE datacenters
(
    id uuid primary key,
    title varchar(255) not null,
    geolocation varchar(255) not null,
    total_size double precision not null default 0,
    created timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE vmachines
(
    id uuid primary key,
    datacenter_id uuid references datacenters(id),
    label varchar(255) not null,
    description varchar(255),
    size int,
    created timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users_vmachines
(
    id uuid primary key,
    user_id uuid references users(id) on delete CASCADE,
    vm_id uuid references vmachines(id) on delete CASCADE
);

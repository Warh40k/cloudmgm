CREATE TABLE users
(
    id uuid primary key,
    name varchar(255) not null,
    username varchar(255) not null unique,
    password_hash varchar(255) not null,
    created timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE vmachines
(
    id uuid primary key,
    title varchar(255) not null,
    description varchar(255),
    status smallint,
    created timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users_vmachines
(
    id uuid primary key,
    user_id uuid references users(id) on delete CASCADE,
    vm_id uuid references vmachines(id) on delete CASCADE
);

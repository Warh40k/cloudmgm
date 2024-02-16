ALTER TABLE users DROP COLUMN login;
ALTER TABLE users ADD COLUMN username varchar(255) not null unique;
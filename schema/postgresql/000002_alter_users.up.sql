ALTER TABLE users ADD COLUMN login varchar(255) not null unique;
ALTER TABLE users DROP COLUMN username;
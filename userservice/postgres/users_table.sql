---- Create users table
CREATE TABLE users ( u jsonb NOT NULL);

-- Default GIN index over entire document
CREATE INDEX idxginusers ON users USING GIN (u);

-- Force username to be unique
CREATE UNIQUE INDEX username ON users( (u->>'username') );
-- Force id to be unique
CREATE UNIQUE INDEX id ON users( (u->>'id') );
-- Force email to be unique
CREATE UNIQUE INDEX email ON users( (u->>'email') );

-- Force username to be non-null
ALTER TABLE users ADD CONSTRAINT username_must_exist CHECK (u ? 'username');
-- Force id to be non-null
ALTER TABLE users ADD CONSTRAINT id_must_exist CHECK (u ? 'id');
-- Force email to be non-null
ALTER TABLE users ADD CONSTRAINT email_must_exist CHECK (u ? 'email');
--Force password to be non-null
ALTER TABLE users ADD CONSTRAINT password_must_exist CHECK (u ? 'password');

-- Force username not to be empty string
ALTER TABLE users ADD CONSTRAINT username_must_not_be_empty_string CHECK ((u->>'username') != '');
-- Force id not to be empty string
ALTER TABLE users ADD CONSTRAINT id_must_not_be_empty_string CHECK ((u->>'id') != '');
-- Force email not to be empty string
ALTER TABLE users ADD CONSTRAINT email_must_not_be_empty_string CHECK ((u->>'email') != '');
-- Force password not to be empty string
ALTER TABLE users ADD CONSTRAINT password_must_not_be_empty_string CHECK ((u->>'password') != '');


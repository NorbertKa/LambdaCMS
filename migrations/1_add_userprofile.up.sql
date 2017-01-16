CREATE TABLE IF NOT EXISTS userprofile (
    id BIGSERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role TEXT DEFAULT 'user' NOT NULL,
    registered timestamp without time zone default (now() at time zone 'utc'),
    lastLogin timestamp without time zone default (now() at time zone 'utc')
);
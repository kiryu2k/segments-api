CREATE TABLE IF NOT EXISTS segment (
    id SERIAL PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS users_segments (
    user_id INTEGER NOT NULL,
    segment_id SERIAL REFERENCES segment (id) ON DELETE CASCADE,
    delete_time TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS logs (
    user_id INTEGER NOT NULL,
    slug TEXT NOT NULL,
    operation TEXT NOT NULL,
    request_time TIMESTAMP NOT NULL
);
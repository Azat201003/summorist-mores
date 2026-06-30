CREATE TABLE IF NOT EXISTS metas (
    more_id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    creator_id INTEGER
);

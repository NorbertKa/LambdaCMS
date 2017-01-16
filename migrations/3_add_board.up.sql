CREATE TABLE IF NOT EXISTS board (
    id BIGSERIAL PRIMARY KEY,
    userId BIGSERIAL references userprofile(id),
    title TEXT NOT NULL UNIQUE,
    miniDescription TEXT NOT NULL,
    fullDescription TEXT NOT NULL,
    image TEXT NOT NULL
)
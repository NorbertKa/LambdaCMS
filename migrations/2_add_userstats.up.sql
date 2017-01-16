CREATE TABLE IF NOT EXISTS userstats (
    id BIGSERIAL PRIMARY KEY,
    userId BIGSERIAL references userprofile(id),
    Posts int DEFAULT 0 NOT NULL,
    Comments int DEFAULT 0 NOT NULL,
    Upvotes int DEFAULT 0 NOT NULL,
    Downvotes int DEFAULT 0 NOT NULL
)
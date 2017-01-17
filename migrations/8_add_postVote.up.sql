CREATE TABLE IF NOT EXISTS postVote (
    id BIGSERIAL PRIMARY KEY,
    userId BIGSERIAL references userprofile(id),
    postId BIGSERIAL references post(id) NOT NULL,
    type TEXT NOT NULL
)
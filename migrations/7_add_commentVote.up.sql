CREATE TABLE IF NOT EXISTS commentVote (
    id BIGSERIAL PRIMARY KEY,
    userId BIGSERIAL references userprofile(id),
    commentId BIGSERIAL references comment(id) NOT NULL,
    type TEXT NOT NULL
)
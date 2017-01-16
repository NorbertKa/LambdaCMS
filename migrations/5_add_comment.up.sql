CREATE TABLE IF NOT EXISTS comment (
    id BIGSERIAL PRIMARY KEY,
    userId BIGSERIAL references userprofile(id),
    postId BIGSERIAL references post(id) NOT NULL,
    commentId BIGSERIAL references comment(id) NOT NULL,
    body TEXT NOT NULL,
    posted timestamp without time zone default (now() at time zone 'utc'),
    upvotes INT DEFAULT 0 NOT NULL,
    downvotes INT DEFAULT 0 NOT NULL
)
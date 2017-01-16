CREATE TABLE IF NOT EXISTS post (
    id BIGSERIAL PRIMARY KEY,
    userId BIGSERIAL references userprofile(id),
    boardId BIGSERIAL references board(id),
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    posted timestamp without time zone default (now() at time zone 'utc'),
    upvotes INT DEFAULT 0 NOT NULL,
    downvotes INT DEFAULT 0 NOT NULL
)
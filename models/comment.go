package db

import "errors"

type Comment struct {
	Id        int    `json:"id"`
	UserId    int    `json:"userId"`
	PostId    int    `json:"postId"`
	CommentId int    `json:"commentId"`
	Body      string `json:"body"`
	Posted    string `json:"posted"`
	Upvotes   int    `json:"upvotes"`
	Downvotes int    `json:"downvotes"`
}

type Comments []Comment

func Comment_GetAll(db *DB) (Comments, error) {
	rows, err := db.Postgre.Query("SELECT id, userId, postId, commentId, body, posted, upvotes, downvotes FROM comment")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comment := Comment{}
	comments := Comments{}
	for rows.Next() {
		err := rows.Scan(&comment.Id, &comment.UserId, &comment.PostId, &comment.CommentId, &comment.Body, &comment.Posted, &comment.Upvotes, &comment.Downvotes)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func Comment_GetByUserId(db *DB, userId int) (Comments, error) {
	rows, err := db.Postgre.Query("SELECT id, userId, postId, commentId, body, posted, upvotes, downvotes FROM comment WHERE userId = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comment := Comment{}
	comments := Comments{}
	for rows.Next() {
		err := rows.Scan(&comment.Id, &comment.UserId, &comment.PostId, &comment.CommentId, &comment.Body, &comment.Posted, &comment.Upvotes, &comment.Downvotes)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func Comment_GetByPostId(db *DB, postId int) (Comments, error) {
	rows, err := db.Postgre.Query("SELECT id, userId, postId, commentId, body, posted, upvotes, downvotes FROM comment WHERE postId = $1", postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comment := Comment{}
	comments := Comments{}
	for rows.Next() {
		err := rows.Scan(&comment.Id, &comment.UserId, &comment.PostId, &comment.CommentId, &comment.Body, &comment.Posted, &comment.Upvotes, &comment.Downvotes)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func Comment_GetByCommentId(db *DB, commentId int) (Comments, error) {
	if commentId == 0 {
		return nil, errors.New("Not a sub-comment")
	}
	rows, err := db.Postgre.Query("SELECT id, userId, postId, commentId, body, posted, upvotes, downvotes FROM comment WHERE commentId = $1", commentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comment := Comment{}
	comments := Comments{}
	for rows.Next() {
		err := rows.Scan(&comment.Id, &comment.UserId, &comment.PostId, &comment.CommentId, &comment.Body, &comment.Posted, &comment.Upvotes, &comment.Downvotes)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func Comment_GetById(db *DB, id int) (Comment, error) {
	rows, err := db.Postgre.Query("SELECT id, userId, postId, commentId, body, posted, upvotes, downvotes FROM comment WHERE id = $1", id)
	if err != nil {
		return Comment{}, err
	}
	defer rows.Close()
	comment := Comment{}
	for rows.Next() {
		err := rows.Scan(&comment.Id, &comment.UserId, &comment.PostId, &comment.CommentId, &comment.Body, &comment.Posted, &comment.Upvotes, &comment.Downvotes)
		if err != nil {
			return Comment{}, err
		}
	}
	err = rows.Err()
	if err != nil {
		return Comment{}, err
	}
	return comment, nil
}

func (c Comment) Create(db *DB) error {
	if len(c.Body) <= 0 {
		return errors.New("Invalid Body")
	}
	user, err := User_GetById(db, c.UserId)
	if err != nil {
		return err
	}
	if user.Id <= 0 {
		errors.New("User not found")
	}
	post, err := Post_GetById(db, c.PostId)
	if err != nil {
		return err
	}
	if post.Id <= 0 {
		errors.New("Post not found")
	}
	if c.CommentId != 0 {
		comment, err := Comment_GetById(db, c.CommentId)
		if err != nil {
			return err
		}
		if comment.Id <= 0 {
			return errors.New("Top comment not found")
		}
	}
	stmt, err := db.Postgre.Prepare("INSERT INTO comment(userId, postId, commentId, body) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(c.UserId, c.PostId, c.CommentId, c.Body)
	if err != nil {
		return err
	}
	return nil
}

func (c Comment) Delete(db *DB) error {
	stmt, err := db.Postgre.Prepare("DELETE FROM comment WHERE id = $1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(c.Id)
	if err != nil {
		return err
	}
	return nil
}

func (c Comment) Update(db *DB) error {
	if len(c.Body) <= 0 {
		return errors.New("Invalid Body")
	}
	stmt, err := db.Postgre.Prepare("UPDATE comment SET body = $1 WHERE id = $2")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(c.Body, c.Id)
	if err != nil {
		return err
	}
	return nil
}

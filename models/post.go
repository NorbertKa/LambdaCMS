package db

import "errors"

type Post struct {
	Id        int    `json:"id"`
	UserId    int    `json:"userId"`
	BoardId   int    `json:"boardId"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Posted    string `json:"posted"`
	Upvotes   int    `json:"upvotes"`
	Downvotes int    `json:"downvotes"`
}

type Posts []Post

func Post_GetAllBody(db *DB, body string) (Posts, error) {
	rows, err := db.Postgre.Query("SELECT id, userId, boardId, title, body, posted, upvotes, downvotes FROM post WHERE body LIKE '%' || $1 || '%' ORDER BY title DESC LIMIT 20", body)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	post := Post{}
	var posts Posts
	for rows.Next() {
		err := rows.Scan(&post.Id, &post.UserId, &post.BoardId, &post.Title, &post.Body, &post.Posted, &post.Upvotes, &post.Downvotes)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func Post_GetAllTitle(db *DB, title string) (Posts, error) {
	rows, err := db.Postgre.Query("SELECT id, userId, boardId, title, body, posted, upvotes, downvotes FROM post WHERE title LIKE '%' || $1 || '%' ORDER BY title DESC LIMIT 20", title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	post := Post{}
	var posts Posts
	for rows.Next() {
		err := rows.Scan(&post.Id, &post.UserId, &post.BoardId, &post.Title, &post.Body, &post.Posted, &post.Upvotes, &post.Downvotes)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func Post_GetAll(db *DB) (Posts, error) {
	rows, err := db.Postgre.Query("SELECT id, userId, boardId, title, body, posted, upvotes, downvotes FROM post")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	post := Post{}
	var posts Posts
	for rows.Next() {
		err := rows.Scan(&post.Id, &post.UserId, &post.BoardId, &post.Title, &post.Body, &post.Posted, &post.Upvotes, &post.Downvotes)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func Post_GetAllByUserId(db *DB, userId int) (Posts, error) {
	rows, err := db.Postgre.Query("SELECT id, userId, boardId, title, body, posted, upvotes, downvotes FROM post WHERE userId = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	post := Post{}
	var posts Posts
	for rows.Next() {
		err := rows.Scan(&post.Id, &post.UserId, &post.BoardId, &post.Title, &post.Body, &post.Posted, &post.Upvotes, &post.Downvotes)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func Post_GetAllByBoardId(db *DB, boardId int) (Posts, error) {
	rows, err := db.Postgre.Query("SELECT id, userId, boardId, title, body, posted, upvotes, downvotes FROM post WHERE boardId = $1 ORDER BY id DESC LIMIT 100", boardId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	post := Post{}
	var posts Posts
	for rows.Next() {
		err := rows.Scan(&post.Id, &post.UserId, &post.BoardId, &post.Title, &post.Body, &post.Posted, &post.Upvotes, &post.Downvotes)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func Post_GetById(db *DB, id int) (Post, error) {
	rows, err := db.Postgre.Query("SELECT id, userId, boardId, title, body, posted, upvotes, downvotes FROM post WHERE id = $1", id)
	if err != nil {
		return Post{}, err
	}
	defer rows.Close()
	post := Post{}
	for rows.Next() {
		err := rows.Scan(&post.Id, &post.UserId, &post.BoardId, &post.Title, &post.Body, &post.Posted, &post.Upvotes, &post.Downvotes)
		if err != nil {
			return Post{}, err
		}
	}
	err = rows.Err()
	if err != nil {
		return Post{}, err
	}
	return post, nil
}

func (p Post) Create(db *DB) error {
	if len(p.Title) <= 0 {
		return errors.New("Invalid title")
	}
	if len(p.Body) <= 0 {
		return errors.New("Invalid Body")
	}
	user, err := User_GetById(db, p.UserId)
	if err != nil {
		return err
	}
	if user.Id <= 0 {
		return errors.New("User not found")
	}
	board, err := Board_GetById(db, p.BoardId)
	if err != nil {
		return err
	}
	if board.Id <= 0 {
		return errors.New("Board not found")
	}
	stmt, err := db.Postgre.Prepare("INSERT INTO post(userId, boardId, title, body) VALUES($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(p.UserId, p.BoardId, p.Title, p.Body)
	if err != nil {
		return err
	}
	return nil
}

func (p Post) Update(db *DB) error {
	if len(p.Title) <= 0 {
		return errors.New("Invalid title")
	}
	if len(p.Body) <= 0 {
		return errors.New("Invalid Body")
	}
	stmt, err := db.Postgre.Prepare("UPDATE post SET title = $1, body = $2 WHERE id = $3")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(p.Title, p.Body, p.Id)
	if err != nil {
		return err
	}
	return nil
}

func (p Post) Delete(db *DB) error {
	stmt, err := db.Postgre.Prepare("DELETE FROM post WHERE id = $1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(p.Id)
	if err != nil {
		return err
	}
	return nil
}

func (p Post) Upvote(db *DB) error {
	stmt, err := db.Postgre.Prepare("UPDATE post SET upvotes = upvotes + 1 WHERE id = $1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(p.Id)
	if err != nil {
		return err
	}
	return nil
}

func (p Post) RemoveUpvote(db *DB) error {
	stmt, err := db.Postgre.Prepare("UPDATE post SET upvotes = upvotes - 1 WHERE id = $1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(p.Id)
	if err != nil {
		return err
	}
	return nil
}

func (p Post) Downvote(db *DB) error {
	stmt, err := db.Postgre.Prepare("UPDATE post SET downvotes = downvotes + 1 WHERE id = $1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(p.Id)
	if err != nil {
		return err
	}
	return nil
}

func (p Post) RemoveDownvote(db *DB) error {
	stmt, err := db.Postgre.Prepare("UPDATE post SET downvotes = downvotes - 1 WHERE id = $1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(p.Id)
	if err != nil {
		return err
	}
	return nil
}

package db

import "errors"

type PostVote struct {
	Id     int    `json:"id"`
	UserId int    `json:"userId"`
	PostId int    `json:"postId"`
	Type   string `json:"type"`
}

func PostVote_GetByUserData(db *DB, userId int, postId int) (PostVote, error) {
	user, err := User_GetById(db, userId)
	if err != nil {
		return PostVote{}, err
	}
	if user.Id == 0 {
		return PostVote{}, errors.New("user not found")
	}
	post, err := Post_GetById(db, postId)
	if err != nil {
		return PostVote{}, err
	}
	if post.Id == 0 {
		return PostVote{}, errors.New("post not found")
	}
	rows, err := db.Postgre.Query("SELECT id, userId, postId, type FROM postVote WHERE userId = $1 AND postId = $2", userId, postId)
	if err != nil {
		return PostVote{}, err
	}
	defer rows.Close()
	postVote := PostVote{}
	for rows.Next() {
		err := rows.Scan(&postVote.Id, &postVote.UserId, &postVote.PostId, &postVote.Type)
		if err != nil {
			return PostVote{}, err
		}
	}
	err = rows.Err()
	if err != nil {
		return PostVote{}, err
	}
	return postVote, nil
}

func (c PostVote) Create(db *DB) error {
	if c.Type == "upvote" || c.Type == "downvote" {

		user, err := User_GetById(db, c.UserId)
		if err != nil {
			return err
		}
		if user.Id == 0 {
			return errors.New("user not found")
		}
		post, err := Post_GetById(db, c.PostId)
		if err != nil {
			return err
		}
		if post.Id == 0 {
			return errors.New("post not found")
		}
		stmt, err := db.Postgre.Prepare("INSERT INTO postVote(userId, postId, type) VALUES($1, $2, $3)")
		if err != nil {
			return err
		}
		_, err = stmt.Exec(c.UserId, c.PostId, c.Type)
		if err != nil {
			return err
		}
		return nil

	} else {
		errors.New("unsupported vote type")
	}
	return nil
}

func (c PostVote) Change(db *DB) error {
	postVote, err := PostVote_GetByUserData(db, c.UserId, c.PostId)
	if err != nil {
		return err
	}
	stmt, err := db.Postgre.Prepare("UPDATE postVote SET type = $1 WHERE id = $2")
	if err != nil {
		return err
	}
	if postVote.Type == "upvote" {
		_, err = stmt.Exec("downvote", c.Id)
		if err != nil {
			return err
		}
	} else if postVote.Type == "downvote" {
		_, err = stmt.Exec("upvote", c.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

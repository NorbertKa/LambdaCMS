package db

import "errors"

type CommentVote struct {
	Id        int    `json:"id"`
	UserId    int    `json:"userId"`
	CommentId int    `json:"commentId"`
	Type      string `json:"type"`
}

func CommentVote_GetByUserData(db *DB, userId int, commentId int) (CommentVote, error) {
	user, err := User_GetById(db, userId)
	if err != nil {
		return CommentVote{}, err
	}
	if user.Id == 0 {
		return CommentVote{}, errors.New("user not found")
	}
	comment, err := Comment_GetById(db, commentId)
	if err != nil {
		return CommentVote{}, err
	}
	if comment.Id == 0 {
		return CommentVote{}, errors.New("comment not found")
	}
	rows, err := db.Postgre.Query("SELECT id, userId, commentId, type FROM commentVote WHERE userId = $1 AND commentId = $2", userId, commentId)
	if err != nil {
		return CommentVote{}, err
	}
	defer rows.Close()
	commentVote := CommentVote{}
	for rows.Next() {
		err := rows.Scan(&commentVote.Id, &commentVote.UserId, &commentVote.CommentId, &commentVote.Type)
		if err != nil {
			return CommentVote{}, err
		}
	}
	err = rows.Err()
	if err != nil {
		return CommentVote{}, err
	}
	return commentVote, nil
}

func (c CommentVote) Create(db *DB) error {
	if c.Type == "upvote" || c.Type == "downvote" {

		user, err := User_GetById(db, c.UserId)
		if err != nil {
			return err
		}
		if user.Id == 0 {
			return errors.New("user not found")
		}
		comment, err := Comment_GetById(db, int(c.CommentId))
		if err != nil {
			return err
		}
		if comment.Id == 0 {
			return errors.New("comment not found")
		}
		stmt, err := db.Postgre.Prepare("INSERT INTO commentVote(userId, commentId, type) VALUES($1, $2, $3)")
		if err != nil {
			return err
		}
		_, err = stmt.Exec(c.UserId, c.CommentId, c.Type)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("Unsupported content type")
	}
}

func (c CommentVote) Change(db *DB) error {
	commentVote, err := CommentVote_GetByUserData(db, c.UserId, c.CommentId)
	if err != nil {
		return err
	}
	stmt, err := db.Postgre.Prepare("UPDATE commentVote SET type = $1 WHERE id = $2")
	if err != nil {
		return err
	}
	if commentVote.Type == "upvote" {
		_, err = stmt.Exec("downvote", c.Id)
		if err != nil {
			return err
		}
	} else if commentVote.Type == "downvote" {
		_, err = stmt.Exec("upvote", c.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

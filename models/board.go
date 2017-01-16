package db

import "errors"

const (
	ErrBoardTitleLength           string = "Board Title should be under 20 characters"
	ErrBoardMiniDescriptionLength string = "Board Title should be under 255 characters"
	ErrBoardFullDescriptionLength string = "Board Title should be under 2000 characters"
)

type Board struct {
	Id              int    `json:"id"`
	UserId          int    `json:"userId"`
	Title           string `json:"title"`
	MiniDescription string `json:"miniDescription"`
	FullDescription string `json:"fullDescription"`
	Image           string `json:"image"`
}

type Boards []Board

func (b Board) Validate() error {
	if len(b.Title) > 20 {
		return errors.New(ErrBoardTitleLength)
	}
	if len(b.MiniDescription) > 255 {
		return errors.New(ErrBoardMiniDescriptionLength)
	}
	if len(b.FullDescription) > 2000 {
		return errors.New(ErrBoardFullDescriptionLength)
	}
	return nil
}

func Board_GetAll(db *DB) (Boards, error) {
	rows, err := db.Postgre.Query("SELECT id, userId, title, miniDescription, fullDescription, image FROM board ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	board := Board{}
	var boards Boards
	for rows.Next() {
		err := rows.Scan(&board.Id, &board.UserId, &board.Title, &board.MiniDescription, &board.FullDescription, &board.Image)
		if err != nil {
			return nil, err
		}
		boards = append(boards, board)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return boards, nil
}

func Board_GetAllLimit(db *DB, limit int) (Boards, error) {
	rows, err := db.Postgre.Query("SELECT id, userId, title, miniDescription, fullDescription, image FROM board ORDER BY id DESC LIMIT $1", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	board := Board{}
	var boards Boards
	for rows.Next() {
		err := rows.Scan(&board.Id, &board.UserId, &board.Title, &board.MiniDescription, &board.FullDescription, &board.Image)
		if err != nil {
			return nil, err
		}
		boards = append(boards, board)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return boards, nil
}

func Board_GetAllTitle(db *DB, title string) (Boards, error) {
	rows, err := db.Postgre.Query("SELECT id, userId, title, miniDescription, fullDescription, image FROM board WHERE title LIKE '%' || $1 || '%'", title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	board := Board{}
	var boards Boards
	for rows.Next() {
		err := rows.Scan(&board.Id, &board.UserId, &board.Title, &board.MiniDescription, &board.FullDescription, &board.Image)
		if err != nil {
			return nil, err
		}
		boards = append(boards, board)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return boards, nil
}

func Board_GetById(db *DB, id int) (*Board, error) {
	rows, err := db.Postgre.Query("SELECT id, userId, title, miniDescription, fullDescription, image FROM board WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	board := Board{}
	for rows.Next() {
		err := rows.Scan(&board.Id, &board.UserId, &board.Title, &board.MiniDescription, &board.FullDescription, &board.Image)
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &board, nil
}

func Board_GetByTitle(db *DB, title string) (*Board, error) {
	rows, err := db.Postgre.Query("SELECT id, userId, title, miniDescription, fullDescription, image FROM board WHERE title = $1", title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	board := Board{}
	for rows.Next() {
		err := rows.Scan(&board.Id, &board.UserId, &board.Title, &board.MiniDescription, &board.FullDescription, &board.Image)
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &board, nil
}

func Board_GetByUserId(db *DB, id int) (Boards, error) {
	rows, err := db.Postgre.Query("SELECT id, userId, title, miniDescription, fullDescription, image FROM board WHERE userId = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	board := Board{}
	var boards Boards
	for rows.Next() {
		err := rows.Scan(&board.Id, &board.UserId, &board.Title, &board.MiniDescription, &board.FullDescription, &board.Image)
		if err != nil {
			return nil, err
		}
		boards = append(boards, board)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return boards, nil
}

func (b Board) Create(db *DB) error {
	err := b.Validate()
	if err != nil {
		return err
	}
	stmt, err := db.Postgre.Prepare("INSERT INTO board(userId, title, miniDescription, fullDescription, image) VALUES($1, $2, $3, $4, $5)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(b.UserId, b.Title, b.MiniDescription, b.FullDescription, b.Image)
	if err != nil {
		return err
	}
	return nil
}

func (b Board) Update(db *DB) error {
	err := b.Validate()
	if err != nil {
		return err
	}
	stmt, err := db.Postgre.Prepare("UPDATE board SET miniDescription = $1, fullDescription = $2, image = $3 WHERE id = $3")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(b.MiniDescription, b.FullDescription, b.Id)
	if err != nil {
		return err
	}
	return nil
}

func (b Board) Delete(db *DB) error {
	stmt, err := db.Postgre.Prepare("DELETE FROM board WHERE id = $1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(b.Id)
	if err != nil {
		return err
	}
	return nil
}

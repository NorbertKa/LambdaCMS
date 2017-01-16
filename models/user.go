package db

import "errors"
import "golang.org/x/crypto/bcrypt"

const (
	ErrUserUsernameLength string = "User username should be under 20 characters"
)

type User struct {
	Id         int    `json:"id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Hash       string `json:"hash"`
	Role       string `json:"role"`
	Registered string `json:"registered"`
	LastLogin  string `json:"lastLogin"`
}

type Users []User

func User_GetAll(db *DB) (Users, error) {
	rows, err := db.Postgre.Query("SELECT id, username, role, registered, lastLogin FROM userprofile")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	user := User{}
	var users Users
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Role, &user.Registered, &user.LastLogin)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func User_GetById(db *DB, id int) (*User, error) {
	rows, err := db.Postgre.Query("SELECT id, username, role, registered, lastLogin FROM userprofile WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	user := User{}
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Role, &user.Registered, &user.LastLogin)
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u User) Validate() error {
	if len(u.Username) > 20 {
		return errors.New(ErrUserUsernameLength)
	}
	return nil
}

func (u User) Create(db *DB) error {
	err := u.Validate()
	if err != nil {
		return err
	}
	stmt, err := db.Postgre.Prepare("INSERT INTO userprofile(username, password, role) VALUES($1, $2, $3)")
	if err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(u.Username, string(hash), u.Role)
	if err != nil {
		return err
	}
	return nil
}

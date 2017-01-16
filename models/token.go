package db

import (
	"errors"
	"strconv"
	"time"
)

type Token string

type Tokens []Token

type UserToken struct {
	Token   Token  `json:"token"`
	User    User   `json:"user"`
	Created string `json:"created"`
}

type UserTokens []UserToken

func Token_GetAll(db *DB) (UserTokens, error) {
	data := db.Redis.LRange("tokens", 0, -1)
	err := data.Err()
	if err != nil {
		return nil, err
	}
	tokenStrings := data.Val()
	userTokens := UserTokens{}
	for _, token := range tokenStrings {
		userData := db.Redis.HMGet(token, "userId", "username", "role", "created")
		err := userData.Err()
		if err != nil {
			return nil, err
		}
		rawData := userData.Val()
		userId, _ := strconv.Atoi(rawData[0].(string))
		userToken := UserToken{
			Token: Token(token),
			User: User{
				Id:       userId,
				Username: rawData[1].(string),
				Role:     rawData[2].(string),
			},
			Created: rawData[3].(string),
		}
		userTokens = append(userTokens, userToken)
	}
	return userTokens, nil
}

func Token_GetAllBanned(db *DB) (Tokens, error) {
	data := db.Redis.LRange("bannedTokens", 0, -1)
	err := data.Err()
	if err != nil {
		return nil, err
	}
	tokenStrings := data.Val()
	tokens := Tokens{}
	for _, token := range tokenStrings {
		tokens = append(tokens, Token(token))
	}
	return tokens, nil
}

func Token_Create(db *DB, user User, token string) error {
	data := db.Redis.RPush("tokens", token)
	err := data.Err()
	if err != nil {
		return err
	}
	dataSet := db.Redis.HMSet(token, map[string]string{
		"userId":   strconv.Itoa(user.Id),
		"username": user.Username,
		"role":     user.Role,
		"created":  strconv.Itoa(int(time.Now().UnixNano())),
	})
	err = dataSet.Err()
	if err != nil {
		return err
	}
	return nil
}

func Token_Ban(db *DB, token string) error {
	tokens, err := Token_GetAll(db)
	if err != nil {
		return err
	}
	check := false
	for _, t := range tokens {
		if string(t.Token) == token {
			check = true
			break
		}
	}
	if !check {
		errors.New("Token not found")
	}
	data := db.Redis.RPush("bannedTokens", token)
	err = data.Err()
	if err != nil {
		return err
	}
	return nil
}

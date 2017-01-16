package db

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/NorbertKa/LambdaCMS/config"
)

type Postgre *sql.DB

func NewPostgre(config *config.Config) (*Postgre, error) {
	dbinfo := fmt.Sprintf("postgres://" + config.Postgre.Username + ":" + config.Postgre.Password + "@" + config.Postgre.Host + ":" + strconv.Itoa(config.Postgre.Port_Int()) + "/" + config.Postgre.Database)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return nil, err
	}
	var postgre Postgre = db
	return &postgre, nil
}

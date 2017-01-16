package db

import "github.com/NorbertKa/LambdaCMS/config"

type DB struct {
	Redis   Redis
	Postgre Postgre
}

func NewDB(conf *config.Config) (*DB, error) {
	redis, err := NewRedis(conf)
	if err != nil {
		return nil, err
	}

	postgre, err := NewPostgre(conf)
	if err != nil {
		return nil, err
	}

	db := DB{
		Redis:   *redis,
		Postgre: *postgre,
	}
	return &db, nil
}

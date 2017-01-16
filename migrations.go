package main

import (
	"strconv"

	"github.com/NorbertKa/LambdaCMS/config"
	_ "github.com/mattes/migrate/driver/postgres"
	"github.com/mattes/migrate/migrate"
)

func MigrateUp(config *config.Config) []error {
	allErrors, ok := migrate.UpSync("postgres://"+config.Postgre.Username+":"+config.Postgre.Password+"@"+config.Postgre.Host+":"+strconv.Itoa(config.Postgre.Port_Int())+"/"+config.Postgre.Database, "./migrations")
	if !ok {
		return allErrors
	}
	return nil
}

func MigrateDown(config *config.Config) []error {
	allErrors, ok := migrate.DownSync("postgres://"+config.Postgre.Username+":"+config.Postgre.Password+"@"+config.Postgre.Host+":"+strconv.Itoa(config.Postgre.Port_Int())+"/"+config.Postgre.Database, "./migrations")
	if !ok {
		return allErrors
	}
	return nil
}

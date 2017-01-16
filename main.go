package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/NorbertKa/LambdaCMS/config"
	"github.com/NorbertKa/LambdaCMS/controllers"
	"github.com/NorbertKa/LambdaCMS/models"
	"github.com/julienschmidt/httprouter"
)

const version string = "0.0.1"

func main() {
	migrationFlag := flag.Bool("m", true, "Migrate database UP or DOWN")
	flag.Parse()
	fmt.Println("[LambdaCMS]")
	fmt.Println("[Version : " + version + "]")
	conf, err := config.ReadConfig("config.json")
	if err != nil {
		fmt.Println("[Error: " + err.Error() + "]")
	}
	fmt.Println("#=> PORT: ", conf.Port_Int())
	fmt.Println("#=> TITLE: ", conf.Title)
	fmt.Println("#=> DESC: ", conf.Description)
	fmt.Println("====== Connecting to DB ======")
	db, err := db.NewDB(conf)
	if err != nil {
		fmt.Println("[Error: " + err.Error() + "]")
	} else {
		fmt.Println("#=> SUCCESS")
	}
	fmt.Println("====== Running Migrations ======")
	if *migrationFlag {
		err := MigrateUp(conf)
		for _, e := range err {
			fmt.Println(e)
		}
	} else {
		err := MigrateDown(conf)
		for _, e := range err {
			fmt.Println(e)
		}
	}
	handler := controller.NewHandler()
	handler.DB = db
	handler.Conf = conf

	router := httprouter.New()
	router.GET("/", handler.Index_GET)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(conf.Port_Int()), router))
}

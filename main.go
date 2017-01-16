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

	router.GET("/me/boards", handler.User_GetBoards)
	router.GET("/users", handler.Users_Get)
	router.POST("/user", handler.User_Post)
	router.GET("/user/:ID", handler.User_Get)
	router.GET("/user/:ID/boards", handler.User_GetBoardsID)
	router.POST("/board", handler.Board_POST)
	router.GET("/board/:ID", handler.Board_GET)
	router.PUT("/board/:ID", handler.Board_UPDATE)
	router.GET("/board/:ID/posts", handler.Board_GET_Posts)
	router.GET("/boards", handler.Boards_GET)
	router.GET("/boards/limit/:LIMIT", handler.Boards_GET_LIMIT)
	router.GET("/boards/search/:TITLE", handler.Boards_GET_TITLE)
	router.POST("/login", handler.Login_POST)
	router.GET("/login", handler.Login_GET)
	router.GET("/tokens", handler.Token_GET)
	router.GET("/posts", handler.Posts_GET)
	router.POST("/post", handler.Post_POST)
	router.GET("/post/:ID", handler.Post_GET)
	router.PUT("/post/:ID", handler.Post_Update)
	router.GET("/admin/getadmin", handler.Admin_GETADMIN)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(conf.Port_Int()), router))
}

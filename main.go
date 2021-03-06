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
	"github.com/rs/cors"
)

const version string = "0.0.1"

func main() {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
	})
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
	router.GET("/me/posts", handler.User_GetPosts)
	router.GET("/me/comments", handler.User_GetComments)
	router.GET("/me", handler.Login_GET)
	router.GET("/users", handler.Users_Get)
	router.POST("/user", handler.User_Post)
	router.GET("/user/:ID", handler.User_Get)
	router.GET("/user/:ID/boards", handler.User_GetBoardsID)
	router.GET("/user/:ID/posts", handler.User_GetPostsID)
	router.GET("/user/:ID/comments", handler.User_GetCommentsID)
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
	router.GET("/post/:ID/upvote", handler.Post_UPVOTE)
	router.GET("/post/:ID/downvote", handler.Post_DOWNVOTE)
	router.GET("/post/:ID/comments", handler.Posts_GetComments)
	router.DELETE("/post/:ID", handler.Post_Delete)
	router.GET("/comments", handler.Comments_GET)
	router.POST("/comment", handler.Comment_POST)
	router.GET("/comment/:ID/comments", handler.Comment_GETComments)
	router.GET("/comment/:ID", handler.Comment_GET)
	router.PUT("/comment/:ID", handler.Comment_UPDATE)
	router.GET("/comment/:ID/upvote", handler.Comment_UPVOTE)
	router.GET("/comment/:ID/downvote", handler.Comment_DOWNVOTE)
	router.DELETE("/comment/:ID", handler.Comment_DELETE)
	router.GET("/admin/getadmin", handler.Admin_GETADMIN)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(conf.Port_Int()), c.Handler(router)))
}

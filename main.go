package main

import (
	"fmt"

	"github.com/NorbertKa/LambdaCMS/config"
)

const version string = "0.0.1"

func main() {
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
}

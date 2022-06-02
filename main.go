package main

import (
	"generatordoc/restapi"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	//fileUrl := "https://sycret.ru/service/apigendoc/forma_025u.xml"
	//err := restapi.DownloadFile("restapi/docs/forma_025u.xml", fileUrl)
	//if err != nil {
	//	panic(err)
	//} else {
	//	fmt.Println("Downloaded: " + fileUrl)
	//	restapi.ParserXml("restapi/docs/forma_025u.xml")
	//}

	chatServerGin := restapi.NewChatServerGin("localhost", 8888)

	chatServerGin.Use()
	chatServerGin.Run()
}

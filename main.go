package main

import (
	"generatordoc/restapi"
)

func main() {

	chatServerGin := restapi.NewChatServerGin("localhost", 8080)
	chatServerGin.Use()
	chatServerGin.Run()
}

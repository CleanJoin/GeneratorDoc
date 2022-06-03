package restapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ChatServerGin struct {
	router          *gin.Engine
	host            string
	port            int
	maxLastMessages uint
}

type GenDoc struct {
	URLTemplate string `json:"URLTemplate"`
	RecordID    int    `json:"RecordID"`
}

type RequestError struct {
	Error string `json:"error"`
}

type IChatServer interface {
	Use()
	Run()
}

func NewChatServerGin(host string, port int) *ChatServerGin {
	chatServerGin := new(ChatServerGin)
	chatServerGin.host = host
	chatServerGin.port = port
	return chatServerGin
}

func (chat *ChatServerGin) Use() {
	chat.router = gin.Default()
	// Конфигурируем все эндпоинты
	chat.router.POST("/gendoc", createGenDocHandler())
	chat.router.GET("/download/:filename", func(ctx *gin.Context) {
		fileName := ctx.Param("filename")
		ctx.FileAttachment("./restapi/docs/"+fileName, fileName)
	})
}

func (chat *ChatServerGin) Run() error {
	if chat.router == nil {
		return fmt.Errorf("gin не сконфигурирован %v", chat.router)
	}
	return chat.router.Run()
}

func createGenDocHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestUser := new(GenDoc)
		var fileDownloadName string
		statusCode, ctx2, checkBadRequest := validateClientRequest(ctx, requestUser)

		if !checkBadRequest {
			ctx.IndentedJSON(statusCode, ctx2)
			return
		}

		err := DownloadXMLFile("restapi/docs/forma_025u.xml", requestUser.URLTemplate)
		if err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Не скачался фаил %s", requestUser.URLTemplate)})
			return
		} else {
			fileDownloadName = AddValueInXML("restapi/docs/forma_025u.xml")
		}
		ctx.IndentedJSON(http.StatusOK, gin.H{"URLWord": fmt.Sprintf("http://localhost:8080/download/%s.doc", fileDownloadName), "URLPDF": fmt.Sprintf("http://localhost:8080/download/%s.pdf", fileDownloadName)})
	}
}

func validateClientRequest(ctx *gin.Context, requestData interface{}) (int, interface{}, bool) {
	err := ctx.BindJSON(&requestData)
	if err != nil {
		return http.StatusBadRequest, gin.H{"error": "Не содержит поля в запросе"}, false
	}
	return http.StatusOK, gin.H{"error": ""}, true
}

package restapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type ChatServerGin struct {
	router          *gin.Engine
	host            string
	port            int
	maxLastMessages uint
}

type ApiGenDoc struct {
	Use      string `json:"Use"`
	Text     string `json:"Text"`
	RecordID string `json:"RecordID"`
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
	chat.router.POST("/apigendoc", docHandler())
	chat.router.POST("/gendoc", createrGenDocHandler())
	chat.router.GET("/health", heathHandler())
}

func (chat *ChatServerGin) Run() error {
	if chat.router == nil {
		return fmt.Errorf("gin не сконфигурирован %v", chat.router)
	}
	return chat.router.Run()
}

func heathHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"success": true,
			"time":    time.Now().Format(time.RFC3339),
		})

	}
}

func docHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestUser := new(ApiGenDoc)

		statusCode, ctx2, checkBadRequest := validateClientRequest(ctx, requestUser)

		if !checkBadRequest {
			ctx.IndentedJSON(statusCode, ctx2)
			return
		}

		jsonValue2, _ := json.Marshal(requestUser)
		responseBody := bytes.NewBuffer(jsonValue2)
		//strings.NewReader(`{"Use":"text_from_use_tag","Text": "data_from_text_tag","RecordID":"30"}`)
		req, err := http.NewRequest("POST", "https://sycret.ru/service/apigendoc/apigendoc", responseBody)
		if err != nil {
			log.Fatal("Error reading request. ", err)
		}
		// Set headers
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		//req.Header.Set("Host", "sycret.ru")
		req.Header.Set("User-Agent", "python-requests/2.27.1")
		// Create and Add cookie to request
		//	cookie := http.Cookie{Name: "cookie_name", Value: "cookie_value"}
		//	req.AddCookie(&cookie)

		// Set client timeout
		client := &http.Client{Timeout: time.Second * 10}

		// Send request
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Error reading response. ", err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Error reading body. ", err)
		}
		satrinout := strings.TrimSpace(string(body))
		//satrinout = strings.Replace(satrinout, `\`, "", -1)

		ctx.IndentedJSON(http.StatusCreated, satrinout)
	}

}

func createrGenDocHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestUser := new(GenDoc)

		statusCode, ctx2, checkBadRequest := validateClientRequest(ctx, requestUser)

		if !checkBadRequest {
			ctx.IndentedJSON(statusCode, ctx2)
			return
		}

		ctx.IndentedJSON(http.StatusOK, gin.H{"URLWord": "your_url\\2022-05-26 14-12-04.doc", "URLPDF": " your_url\\2022-05-26 14-12-04.pdf"})
	}

}

func validateClientRequest(ctx *gin.Context, requestData interface{}) (int, interface{}, bool) {

	err := ctx.BindJSON(&requestData)
	if err != nil {

		return http.StatusBadRequest, gin.H{"error": "Не содержит поля в запросе"}, false
	}
	return http.StatusOK, gin.H{"error": ""}, true
}

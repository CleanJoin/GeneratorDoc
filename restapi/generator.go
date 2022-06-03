package restapi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type ApiGenDoc struct {
	Use      string `json:"Use"`
	Text     string `json:"Text"`
	RecordID string `json:"RecordID"`
}

type ResultApiGenDoc struct {
	Result      string `json:"result"`
	Description string `json:"resultdescription"`
	Data        string `json:"resultdata"`
}

func (apiGenDoc *ApiGenDoc) PostAPIGenDoc(url, use, text, recordID string) *ResultApiGenDoc {
	apiGenDoc.Use = use
	apiGenDoc.Text = text
	apiGenDoc.RecordID = recordID
	resultApiGenDoc := new(ResultApiGenDoc)
	jsonValue2, _ := json.Marshal(apiGenDoc)
	responseBody := bytes.NewBuffer(jsonValue2)
	req, err := http.NewRequest("POST", url, responseBody)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "python-requests/2.27.1")

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
	json.Unmarshal(body, resultApiGenDoc)

	return resultApiGenDoc
}

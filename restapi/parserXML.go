package restapi

import (
	"fmt"
	"github.com/ConvertAPI/convertapi-go"
	"github.com/ConvertAPI/convertapi-go/config"
	"github.com/ConvertAPI/convertapi-go/param"
	"github.com/beevik/etree"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func DownloadXMLFile(filepath string, url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "python-requests/2.27.1")

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	return err
}

func AddValueInXML(filepath string) string {

	doc := etree.NewDocument()
	if err := doc.ReadFromFile(filepath); err != nil {
		panic(err)
	}

	root := doc.SelectElement("wordDocument")
	body := root.SelectElement("body")
	sect := body.SelectElement("sect")
	use := sect.SelectElements("use")
	useTbl := use[0].SelectElements("tbl")
	tr8 := useTbl[3].SelectElement("tr")
	tr10 := useTbl[4].SelectElement("tr")
	tcALL8 := tr8.SelectElements("tc")
	p8 := tcALL8[1].SelectElement("p")
	cardNumber := p8.SelectElement("text")
	tcALL10 := tr10.SelectElements("tc")
	p10 := tcALL10[1].SelectElement("p")
	useClient := p10.SelectElements("use")
	useClientText := useClient[0].SelectElements("text")

	apiGenDoc := new(ApiGenDoc)
	reqApiDocGen := apiGenDoc.PostAPIGenDoc("https://sycret.ru/service/apigendoc/apigendoc", use[0].Attr[2].Value, cardNumber.Attr[0].Value, "30")
	cardNumberR := cardNumber.SelectElement("r")
	cardNumberText := cardNumberR.SelectElement("t")
	cardNumberText.SetText(reqApiDocGen.Data)

	for i, val := range useClientText {
		fmt.Println(val, i)
		reqApiDocGen := apiGenDoc.PostAPIGenDoc("https://sycret.ru/service/apigendoc/apigendoc", useClient[0].Attr[0].Value, val.Attr[0].Value, "30")
		client := useClientText[i].SelectElement("r")
		client = client.SelectElement("t")
		client.SetText(reqApiDocGen.Data)
	}

	t := time.Now()
	doc.WriteToFile(fmt.Sprintf("restapi/docs/%v.doc", t.Format("2006-02-01 15-04-05")))

	config.Default.Secret = "ynP4oWlnrMEVPBAZ"
	convertapi.ConvDef("doc", "pdf",
		param.NewPath("File", fmt.Sprintf("restapi/docs/%v.doc", t.Format("2006-02-01 15-04-05")), nil),
	).ToPath(fmt.Sprintf("restapi/docs/%v.pdf", t.Format("2006-02-01 15-04-05")))
	return t.Format("2006-02-01 15-04-05")
}

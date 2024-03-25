package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/joho/godotenv"
	"io"
	"os"
	"time"
)

type Auditlog struct {
	Time     string
	Url      string
	Method   string
	Request  any
	Response any
}

func getESClient() (*elasticsearch.Client, error) {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	elasticUrl := os.Getenv("ELASTIC_URL")

	cfg := elasticsearch.Config{
		Addresses: []string{elasticUrl},
	}
	es, err := elasticsearch.NewClient(cfg)

	if err != nil {
		return nil, err
	}

	return es, nil
}

func (a *Auditlog) StoreToES() {
	es, err := getESClient()
	if err != nil {
		fmt.Println(err.Error())
	}

	t := time.Now()
	formattedTime := t.UTC().Format("2006-01-02T15:04:05.999999Z")
	a.Time = formattedTime

	body, err := json.Marshal(a)
	if err != nil {
		fmt.Println(err.Error())
	}

	//fmt.Println("Body JSON", string(body))

	// Index the document into Elasticsearch
	res, err := es.Index(
		"auditlogs3",
		bytes.NewReader(body),
		es.Index.WithDocumentID(fmt.Sprintf("%d", t.UnixNano())),
	)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	if res.IsError() {
		fmt.Println(res.Status())
	}

	fmt.Println("Document indexed successfully.")
	//fmt.Println("Time:", now)
	//fmt.Println("Method:", a.Method)
	//fmt.Println("Url:", a.Url)
	//fmt.Println("Request:", a.Request)
	//fmt.Println("Response:", a.Response)
}

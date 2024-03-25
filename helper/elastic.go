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

func (a *Auditlog) StoreToES() error {
	err := godotenv.Load(".env")

	if err != nil {
		return err
	}

	elasticIndex := os.Getenv("ELASTIC_INDEX")

	es, err := getESClient()
	if err != nil {
		return err
	}

	t := time.Now()
	formattedTime := t.UTC().Format("2006-01-02T15:04:05.999999Z")
	a.Time = formattedTime

	body, err := json.Marshal(a)
	if err != nil {
		return err
	}

	res, err := es.Index(
		elasticIndex,
		bytes.NewReader(body),
		es.Index.WithDocumentID(fmt.Sprintf("%d", t.UnixNano())),
	)

	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	if res.IsError() {
		return err
	}
	return nil
}

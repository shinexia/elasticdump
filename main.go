package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/estransport"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

type Opt struct {
	filename  string
	url       string
	host      string
	index     string
	opType    string
	batchSize int
	max       int
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "file, f",
			Usage:       "file",
			Required:    true,
			Value:       "",
			DefaultText: "",
		},
		&cli.StringFlag{
			Name:        "url",
			Usage:       "ElasticSearch url",
			Required:    true,
			Value:       "",
			DefaultText: "",
		},
		&cli.StringFlag{
			Name:        "type",
			Usage:       "load or dump type: [data, mapping]",
			Required:    false,
			Value:       "data",
			DefaultText: "data",
		},
		&cli.IntFlag{
			Name:        "batch",
			Usage:       "batch size",
			Required:    false,
			Value:       1000,
			DefaultText: "1000",
		},
		&cli.IntFlag{
			Name:        "max",
			Usage:       "max records",
			Required:    false,
			Value:       0,
			DefaultText: "0",
		},
	}
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "load",
				Usage:  "load records from file",
				Action: load,
				Flags:  flags,
			},
			{
				Name:   "dump",
				Usage:  "dump records to file",
				Action: dump,
				Flags:  flags,
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func parseOpt(c *cli.Context) *Opt {
	opt := &Opt{
		filename:  c.String("file"),
		url:       c.String("url"),
		opType:    c.String("type"),
		batchSize: c.Int("batch"),
		max:       c.Int("max"),
		index:     "",
		host:      "",
	}
	up, err := url.Parse(opt.url)
	if err != nil {
		log.Fatalf("parse url err: %v\n", err)
	}
	opt.host = up.Scheme + "://" + up.Host
	opt.index = strings.TrimLeft(up.Path, "/")
	if strings.Contains(opt.index, "/") {
		log.Fatalf("invalid index: %v\n", opt.url)
	}
	log.Printf("opt: %#v\n", opt)
	return opt
}

func dump(c *cli.Context) error {
	opt := parseOpt(c)
	switch opt.opType {
	case "data":
		return dumpData(opt)
	case "mapping":
		return dumpMapping(opt)
	}
	log.Fatalf("unknown type: %v\n", opt.opType)
	return nil
}

func load(c *cli.Context) error {
	opt := parseOpt(c)
	switch opt.opType {
	case "data":
		return loadData(opt)
	case "mapping":
		return loadMapping(opt)
	}
	return fmt.Errorf("unknown type: %v\n", opt.opType)
}

func dumpMapping(opt *Opt) error {
	conf := elasticsearch.Config{
		Addresses:         []string{opt.url},
		EnableMetrics:     true,
		EnableDebugLogger: true,
		Logger: &estransport.TextLogger{
			Output:             os.Stdout,
			EnableRequestBody:  true,
			EnableResponseBody: false,
		},
	}
	es, err := elasticsearch.NewClient(conf)
	if err != nil {
		log.Fatalf("create es client err: %v\n", err)
	}
	var (
		startTime = time.Now()
		writer    = newLazyWriter(opt.filename)
	)
	defer writer.Close()
	res, err := es.Info()
	if err != nil {
		log.Fatalf("info err: %v\n", err)
	}
	// If the whole request failed, print error and mark all documents as failed
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("read res body err: %v\n", err)
	}
	res.Body.Close()
	nbytes, err := writer.Write(resBody)
	if err != nil {
		log.Fatalf("write mapping err: %v\n", err)
	}
	cost := time.Now().Sub(startTime).Milliseconds()
	log.Printf("dump mapping cost: %.3fs, nbytes: %v\n", float64(cost)*0.001, nbytes)
	return nil
}

func loadMapping(opt *Opt) error {
	conf := elasticsearch.Config{
		Addresses:         []string{opt.host},
		EnableMetrics:     true,
		EnableDebugLogger: true,
		Logger: &estransport.TextLogger{
			Output:             os.Stdout,
			EnableRequestBody:  true,
			EnableResponseBody: false,
		},
	}
	es, err := elasticsearch.NewClient(conf)
	if err != nil {
		log.Fatalf("create es client err: %v\n", err)
	}
	file, err := os.Open(opt.filename)
	if err != nil {
		log.Fatalf("read mapping failed, err: %v\n", err)
	}
	var raw map[string]json.RawMessage
	err = json.NewDecoder(file).Decode(&raw)
	if err != nil {
		log.Fatalf("parse mapping failed, err: %v\n", err)
	}
	if len(raw) != 1 {
		log.Fatalf("multi index count: %v\n", len(raw))
	}
	var root json.RawMessage
	for _, v := range raw {
		root = v
	}
	var rootMap map[string]json.RawMessage
	err = json.NewDecoder(bytes.NewReader(root)).Decode(&rootMap)
	if err != nil {
		log.Fatalf("parse rootMap failed, err: %v\n", err)
	}
	mappings := rootMap["mappings"]
	settings := rootMap["settings"]
	var settingsMap map[string]json.RawMessage
	err = json.NewDecoder(bytes.NewReader(settings)).Decode(&settingsMap)
	if err != nil {
		log.Fatalf("parse settingsMap failed, err: %v\n", err)
	}
	index := settingsMap["index"]
	var indexMap map[string]json.RawMessage
	err = json.NewDecoder(bytes.NewReader(index)).Decode(&indexMap)
	if err != nil {
		log.Fatalf("parse indexMap failed, err: %v\n", err)
	}
	for _, key := range []string{"creation_date", "uuid", "version", "provided_name"} {
		if _, ok := indexMap[key]; ok {
			delete(indexMap, key)
		}
	}
	index, err = json.Marshal(indexMap)
	if err != nil {
		log.Fatalf("marshal indexMap failed, err: %v\n", err)
	}
	settingsMap["index"] = index
	log.Printf("indexMap: %v\n", string(index))
	req := make(map[string]json.RawMessage, 2)
	req["mappings"] = mappings
	req["settings"], err = json.Marshal(settingsMap)
	if err != nil {
		log.Fatalf("marshal settingsMap failed, err: %v\n", err)
	}
	reqData, err := json.Marshal(req)
	if err != nil {
		log.Fatalf("marshal req failed, err: %v\n", err)
	}
	var (
		startTime = time.Now()
	)
	res, err := es.Indices.Delete([]string{opt.index})
	if err != nil {
		log.Fatalf("delete index err: %v\n", err)
	}
	log.Printf("delete res: %v\n", res)
	res.Body.Close()
	res, err = es.Indices.Create(opt.index, es.Indices.Create.WithBody(bytes.NewReader(reqData)))
	if err != nil {
		log.Fatalf("create index err: %v\n", err)
	}
	log.Printf("create res: %v\n", res)
	res.Body.Close()
	cost := time.Now().Sub(startTime).Milliseconds()
	log.Printf("load mapping cost: %.3fs\n", float64(cost)*0.001)
	return nil
}

func dumpData(opt *Opt) error {
	conf := elasticsearch.Config{
		Addresses:         []string{opt.url},
		EnableMetrics:     true,
		EnableDebugLogger: true,
		Logger: &estransport.TextLogger{
			Output:             os.Stdout,
			EnableRequestBody:  false,
			EnableResponseBody: false,
		},
	}
	es, err := elasticsearch.NewClient(conf)
	if err != nil {
		log.Fatalf("create es client err: %v\n", err)
	}
	var (
		from      = 0
		numDumped = 0
		startTime = time.Now()
		writer    = newLazyWriter(opt.filename)
		scrollId  = ""
	)
	defer writer.Close()
	for {
		var (
			sr  = ScrollResponse{}
			res *esapi.Response
			err error
		)
		if scrollId == "" {
			res, err = es.Search(es.Search.WithContext(context.Background()), es.Search.WithIndex(opt.index), es.Search.WithScroll(time.Minute), es.Search.WithSize(opt.batchSize), es.Search.WithQuery("*:*"))
			if err != nil {
				log.Fatalf("search err: %v\n", err)
			}
		} else {
			req := []byte(fmt.Sprintf(`{"scroll": "%s","scroll_id": "%s"}`, "1m", scrollId))
			log.Printf("req: %v\n", string(req))
			res, err = es.Scroll(es.Scroll.WithContext(context.Background()), es.Scroll.WithBody(bytes.NewReader(req)))
			if err != nil {
				log.Fatalf("search err: %v\n", err)
			}
		}
		// If the whole request failed, print error and mark all documents as failed
		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatalf("read res body err: %v\n", err)
		}
		res.Body.Close()
		if res.IsError() {
			log.Fatalf("Error [%d]: %s\n", res.StatusCode, string(resBody))
		} else {
			if err := json.Unmarshal(resBody, &sr); err != nil {
				log.Fatalf("Failure to to parse response body: %s", err)
			} else {
				scrollId = sr.ScrollID
				succeed := len(sr.Hits.Hits)
				if succeed > 0 {
					nbytes := 0
					for _, record := range sr.Hits.Hits {
						n, err := writer.Write(record)
						if err != nil {
							log.Fatalf("write record failed, err: %v\n", err)
						}
						nbytes += n
						n, err = writer.Write([]byte("\n"))
						if err != nil {
							log.Fatalf("write record failed, err: %v\n", err)
						}
						nbytes += n
					}
					numDumped += succeed
					from += succeed
					log.Printf("indexed succeed: %v, nbytes: %v\n", succeed, nbytes)
				} else {
					log.Printf("res: %v, break\n", string(resBody))
					break
				}
			}
		}
	}
	cost := time.Now().Sub(startTime).Milliseconds()
	log.Printf("num dumped: %v, cost: %.3fs\n", numDumped, float64(cost)*0.001)
	return nil
}

func loadData(opt *Opt) error {
	conf := elasticsearch.Config{
		Addresses:         []string{opt.url},
		EnableMetrics:     true,
		EnableDebugLogger: true,
		Logger: &estransport.TextLogger{
			Output:             os.Stdout,
			EnableRequestBody:  false,
			EnableResponseBody: false,
		},
	}
	es, err := elasticsearch.NewClient(conf)
	if err != nil {
		log.Fatalf("create es client err: %v\n", err)
	}
	records, err := loadRecords(opt.filename)
	if err != nil {
		log.Fatalf("load records err: %v\n", err)
	}
	log.Printf("records: %v\n", len(records))
	if opt.max > 0 {
		records = records[:opt.max]
		log.Printf("sub records: %v\n", len(records))
	}
	var (
		numErrors  int
		numIndexed int
	)
	for i := 0; i < len(records); i += opt.batchSize {
		var (
			buf bytes.Buffer
			blk BulkResponse
		)
		j := i + opt.batchSize
		if j > len(records) {
			j = len(records)
		}
		sub := records[i:j]
		for _, r := range sub {
			meta := []byte(fmt.Sprintf(`{"create": {"_index": "%s", "_type": "_doc", "_id": "%s"}}%s`, opt.index, r.ID, "\n"))
			buf.Write(meta)
			buf.Write(r.Source)
			buf.Write([]byte("\n"))
		}
		data := buf.Bytes()
		res, err := es.Bulk(bytes.NewReader(data))
		if err != nil {
			log.Fatalf("bulk err: %v\n", err)
		}
		// If the whole request failed, print error and mark all documents as failed
		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatalf("read res body err: %v\n", err)
		}
		res.Body.Close()
		if res.IsError() {
			numErrors += j - i
			log.Printf("Error [%d]: %s\n", res.StatusCode, string(resBody))
		} else {
			if err := json.Unmarshal(resBody, &blk); err != nil {
				log.Fatalf("Failure to to parse response body: %s", err)
			} else {
				succeed := 0
				for _, d := range blk.Items {
					// ... so for any HTTP status above 201 ...
					//
					if d.Create.Status > 201 {
						numErrors++
						log.Printf("Error [%d]: %s\n", d.Create.Status, d.Create.Error)
					} else {
						succeed++
					}
				}
				numIndexed += succeed
				log.Printf("indexed succeed: %v\n", succeed)
			}
		}
	}
	log.Printf("num indexed: %v, num errors: %v\n", numIndexed, numErrors)
	return nil
}

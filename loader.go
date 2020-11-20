package main

import (
	"bufio"
	"encoding/json"
	"os"
)

func loadRecords(filename string) ([]*SourceWrap, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	bufSize := 1024 * 1024 * 1024
	scanner.Buffer(make([]byte, bufSize), bufSize)
	scanner.Split(bufio.ScanLines)
	var records []*SourceWrap
	for scanner.Scan() {
		line := scanner.Bytes()
		var src = &Source{}
		err := json.Unmarshal(line, src)
		if err != nil {
			return nil, err
		}
		var rec = &SourceWrap{
			Source: src,
			Data:   json.RawMessage(line),
		}
		records = append(records, rec)
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	return records, nil
}

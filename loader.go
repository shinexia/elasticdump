package main

import (
	"bufio"
	"encoding/json"
	"os"
)

func loadRecords(filename string) ([]*DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var records []*DataRecord
	for scanner.Scan() {
		line := scanner.Bytes()
		var rec = &DataRecord{}
		err := json.Unmarshal(line, rec)
		if err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	return records, nil
}

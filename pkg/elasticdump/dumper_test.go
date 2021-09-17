/*
Copyright 2021 Shine Xia <shine.xgh@gmail.com>.

Licensed under the MIT License.
*/

package elasticdump

import (
	"fmt"
	"testing"
	"time"
)

func newTestDumper() (*Dumper, error) {
	client, err := newTestClient()
	if err != nil {
		return nil, err
	}
	return NewDumper(client), nil
}

func TestDumpperGenTestData(t *testing.T) {
	dumper, err := newTestDumper()
	if err != nil {
		t.Error(err)
		return
	}
	err = dumper.DeleteIndex(testIndexName)
	if err != nil {
		t.Error(err)
		return
	}
	err = dumper.GenTestData(testIndexName, 10, 100)
	if err != nil {
		fmt.Printf("%+v\n", err)
		t.Error(err)
		return
	}
	fmt.Printf("succeed\n")
}

func TestDumpperDumpData(t *testing.T) {
	dumper, err := newTestDumper()
	if err != nil {
		t.Error(err)
		return
	}
	err = dumper.DumpData(testIndexName, "/tmp/"+testIndexName+"-data.json", 100, -1, time.Minute)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("succeed\n")
}

func TestDumpperLoadData(t *testing.T) {
	dumper, err := newTestDumper()
	if err != nil {
		t.Error(err)
		return
	}
	index := testIndexName + "2"
	err = dumper.DeleteIndex(index)
	if err != nil {
		t.Error(err)
		return
	}
	err = dumper.LoadData(index, "/tmp/"+testIndexName+"-data.json", 100, -1, 1024*1024)
	if err != nil {
		fmt.Printf("%+v\n", err)
		t.Error(err)
		return
	}
	fmt.Printf("succeed\n")
}

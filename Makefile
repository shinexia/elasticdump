
.PHONY: clean

elasticsearch-dump: $(wildcard *.go)
	go build -v -o $@

clean:
	rm -fr esdump elasticsearch-dump


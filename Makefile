
.PHONY: clean

esdump: $(wildcard *.go)
	go build -v -o $@

clean:
	rm -fr esdump


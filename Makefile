
.PHONY: clean

elasticdump: $(wildcard *.go)
	go build -v -o $@

clean:
	rm -fr elasticdump


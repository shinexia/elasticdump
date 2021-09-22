# elasticdump

a simple elasticsearch dump & load tool. 

## Why was elasticdump created?

The ElasicSearch did not offer a dump tool, and the only tool provided at <https://github.com/elasticsearch-dump/elasticsearch-dump> depends on a nodejs environment, which is not convenient.

## INSTALL

1. `go get https://github.com/shinexia/elasticdump`

2. or download a prebuilt binary here: <https://github.com/shinexia/elasticdump/releases/>

3. or build from source

``` bash
git clone https://github.com/shinexia/elasticdump.git
cd elasticdump
make
```

## EXAMPLE

``` bash
$ ./elasticdump 


    ┌──────────────────────────────────────────────────────────┐
    │ ElasticDump                                              │
    │ a simple dump/load data/mapping from elasticsearch       │
    │                                                          │
    │ Please give us feedback at:                              │
    │ https://github.com/shinexia/elasticdump/issues           │
    └──────────────────────────────────────────────────────────┘

Example usage:

        elasticdump --host http://localhost:9200 --index elasticdumptest gen  testdata -v=10

        elasticdump --host http://localhost:9200 --index elasticdumptest dump mapping

        elasticdump --host http://localhost:9200 --index elasticdumptest dump data -v=4

        elasticdump --host http://localhost:9200 --index elasticdumptest load mapping --delete

        elasticdump --host http://localhost:9200 --index elasticdumptest load data

Usage:
  elasticdump [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  delete      delete index from elasticsearch
  dump        dump mapping/data from elasticsearch
  gen         gen testdata to elasticsearch
  help        Help about any command
  load        load mapping/data to elasticsearch

Flags:
  -h, --help   help for elasticdump

Use "elasticdump [command] --help" for more information about a command.

```

## LICENSE

MIT License

# elasticdump

a simple elasticsearch dump & load tool. 

## Description

The Elasicsearch did not offer a dump tool, and the only tool provided at <https://github.com/elasticsearch-dump/elasticsearch-dump> is too slowly, and depend on a nodejs environment.

## Install

``` bash
git clone https://github.com/shinexia/elasticdump.git
go build -v
```

or through prebuild binaries

``` bash
curl -fsSL https://github.com/shinexia/elasticdump/releases/download/v0.2.0/elasticdump-linux-x86_64 -o elasticdump
chmod a+x elasticdump
```

## Usage

``` bash
$./elasticdump --help
NAME:
   elasticdump - A new cli application

USAGE:
   elasticdump [global options] command [command options] [arguments...]

COMMANDS:
   load     load records from file
   dump     dump records to file
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)

```

### dump mappings

``` bash
./elasticdump dump --url http://<user_name>:<url_encoded_password>@localhost:9200/myindex --file /tmp/myindex-mapping.json --type mapping
```

### load mappings

``` bash
./elasticdump load --url http://localhost:9200/myindex --file /tmp/myindex-mapping.json --type mapping
```

### dump data

``` bash
./elasticdump dump --url http://localhost:9200/myindex --file /tmp/myindex-data.json  --type data
```

### load data

``` bash
./elasticdump load --url http://localhost:9200/myindex  --file /tmp/myindex-data.json--type data
```

## LICENSE

Apache License 2.0

# elasticsearch-dump

a simple elasticsearch dump & load tool

## Build

``` bash
git clone https://github.com/shinexia/elasticsearch-dump.git
go build -v
```

## Usage

``` bash
$./elasticsearch-dump --help
NAME:
   elasticsearch-dump - A new cli application

USAGE:
   elasticsearch-dump [global options] command [command options] [arguments...]

COMMANDS:
   load     load records from file
   dump     dump records to file
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)

```

### dump mappings

``` bash
./elasticsearch-dump dump --file /tmp/myindex-mapping.json --url http://localhost:9200/myindex --type mapping
```

### load mappings

``` bash
./elasticsearch-dump load --file /tmp/myindex-mapping.json --url http://localhost:9200/myindex --type mapping
```

### dump data

``` bash
./elasticsearch-dump dump --file /tmp/myindex-data.json --url http://localhost:9200/myindex --type data
```

### load data

``` bash
./elasticsearch-dump load --file /tmp/myindex-data.json --url http://localhost:9200/myindex --type data
```

## LICENSE

Apache License 2.0

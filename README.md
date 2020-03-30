# esdump

a simple elasticsearch dump & load tool

## BUILD

``` bash
git clone https://github.com/shinexia/esdump.git
go build -v
```

## Usage

``` bash
$./esdump --help
NAME:
   esdump - A new cli application

USAGE:
   esdump [global options] command [command options] [arguments...]

COMMANDS:
   load     load records from file
   dump     dump records to file
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)

```

### dump mappings

``` bash
./esdump dump --file /tmp/myindex-mapping.json --url http://localhost:9200/myindex --type mapping
```

### load mappings

``` bash
./esdump load --file /tmp/myindex-mapping.json --url http://localhost:9200/myindex --type mapping
```

### dump data

``` bash
./esdump dump --file /tmp/myindex-data.json --url http://localhost:9200/myindex --type data
```

### load data

``` bash
./esdump load --file /tmp/myindex-data.json --url http://localhost:9200/myindex --type data
```

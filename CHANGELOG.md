# CHANGELOGs

## v0.3.8

1. add ldflags to trim binary size

## v0.3.7

1. refactor cmd code structrue

## v0.3.6

2021.09.16

1. reorganize code structure: move importable codes to `pkg/elasticdump` and cmd codes to `cmd/elasticdump`
2. data records format changed: include complete hit record `{"_index":"","_type":"_doc","_id":"","_score":1.0,"_source": {}}`, and previous version only save `_source` field data
3. cli arguments format changed, run `elasticdump -h` to see details
4. add `k8s.io/klog` to support leveled logging, which can configured by `-v=<4-10>` when running cmd
5. optimized dump and load speed, split read/write action to different goroutine


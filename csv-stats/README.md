# CSV Stats

## Benchmarking

In order to benchmark the program, you need to:
- Download [source code](https://media.pragprog.com/titles/rggo/code/rggo-code.zip) for book.
- Copy `performance/colStatsBenchmarkData.tar.gz` from source code to `./testdata` in csv-stats.
- Run following commands to prepare data
```bash
$ cd testdata
$ tar -xzvf colStatsBenchmarkData.tar.gz
```
- Run benchmarking
```bash
$ go test -bench . -run ^$
```
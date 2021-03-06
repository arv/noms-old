# CSV Importer

Imports a CSV file as `List<T>` where `T` is a struct with fields corresponding to the CSV's column headers. The struct spec can also be set manually with the `-header` flag.

## Usage

```
$ cd importer
$ go build
$ ./importer --h=http://localhost:8000  --ds=foo <PATH>
```

## Some places for CSV files

- https://data.cityofnewyork.us/api/views/kku6-nxdu/rows.csv?accessType=DOWNLOAD
- http://www.opendatacache.com/

# CSV Exporter

Export a dataset in CSV format to stdout with column headers.

## Usage

```
$ cd exporter
$ go build
$ ./exporter --h=http://localhost:8000  --ds=foo
```

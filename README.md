[![Build Status](https://travis-ci.org/kshedden/datareader.svg?branch=master)](https://travis-ci.org/kshedden/datareader)
[![Go Report Card](https://goreportcard.com/badge/github.com/kshedden/datareader)](https://goreportcard.com/report/github.com/kshedden/datareader)
[![codecov](https://codecov.io/gh/kshedden/datareader/branch/master/graph/badge.svg)](https://codecov.io/gh/kshedden/datareader)
[![GoDoc](https://godoc.org/github.com/kshedden/datareader?status.png)](https://godoc.org/github.com/kshedden/datareader)

datareader : read SAS and Stata files in Go
=========================

__datareader__ is a pure [Go](https://golang.org) (Golang) package
that can read binary SAS format (SAS7BDAT) and Stata format (dta) data
files into native Go data structures.  For non-Go users, there are
command line utilities that convert SAS and Stata files into text/csv
and parquet files.

The Stata reader is based on the Stata documentation for the [dta file
format](http://www.stata.com/help.cgi?dta) and supports dta versions
115, 117, and 118.

There is no official documentation for SAS binary format files.  The
code here is translated from the Python
[sas7bdat](https://pypi.python.org/pypi/sas7bdat) package, which in
turn is based on an [R
package](https://github.com/BioStatMatt/sas7bdat).  Also see
[here](https://cran.r-project.org/web/packages/sas7bdat/vignettes/sas7bdat.pdf)
for more information about the SAS7BDAT file structure.

This package also provides a simple column-oriented data container
called a `Series`.  Both the SAS reader and Stata reader return the
data as an array of `Series` objects, corresponding to the columns of
the data file.  These can in turn be converted to other formats as
needed.

Both the Stata and SAS reader support streaming access to the data
(i.e. reading the file by chunks of consecutive records).

## SAS

Here is an example of how the SAS reader can be used in a Go program
(error handling omitted for brevity):

```
import (
        "datareader"
        "os"
)

// Create a SAS7BDAT object
f, _ := os.Open("filename.sas7bdat")
sas, _ := datareader.NewSAS7BDATReader(f)

// Read the first 10000 records (rows)
ds, _ := sas.Read(10000)

// If column 0 contains numeric data
// x is a []float64 containing the dta
// m is a []bool containing missingness indicators
x, m, _ := ds[0].AsFloat64Slice()

// If column 1 contains text data
// x is a []string containing the dta
// m is a []bool containing missingness indicators
x, m, _ := ds[1].AsStringSlice()
```

## Stata

Here is an example of how the Stata reader can be used in a Go program
(again with no error handling):

```
import (
        "datareader"
        "os"
)

// Create a StataReader object
f,_ := os.Open("filename.dta")
stata, _ := datareader.NewStataReader(f)

// Read the first 10000 records (rows)
ds, _ := stata.Read(10000)
```

## CSV

The package includes a CSV reader with type inference for the column data types.

```
import (
        "datareader"
)

f, _ := os.Open("filename.csv")
rt := datareader.NewCSVReader(f)
rt.HasHeader = true
dt, _ := rt.Read(-1)
// obtain data from dt as in the SAS example above
```

## Command line utilities

We provide two command-line utilities allowing conversion of SAS and
Stata datasets to other formats without using Go directly.
Executables for several OS's and architectures are contained in the
`bin` directory.  The script used to cross-compile these binaries is
`build.sh`.  To build and install the commands for your local
architecture only, run the Makefile (the executables will be copied
into your GOBIN directory).

The `stattocsv` command converts a SAS7BDAT or Stata dta file to a csv
file, it can be used as follows:

```
> stattocsv file.sas7bdat > file.csv
> stattocsv file.dta > file.csv
```

The `columnize` command takes the data from either a SAS7BDAT or a
Stata dta file, and writes the data from each column into a separate
file.  Numeric data can be stored in either binary (native 8 byte
floats) or text format (binary is considerably faster).

```
> columnize -in=file.sas7bdat -out=cols -mode=binary
> columnize -in=file.dta -out=cols -mode=text
```

## Parquet conversion

We provide a simple and efficient way to convert a SAS7BDAT file to
parquet format, using the
[parquet-go](https://github.com/xitongsys/parquet-go) package.  To
convert a SAS file called 'mydata.sas7bdat' to Parquet format, begin
by running sas_to_parquet as follows:

```
sas_to_parquet -sasfile=mydata.sas7bdat -outdir=. -structname=MyStruct -pkgname=mypackage
```

If you want the Parquet file for use outside of Go, you can specify
any values for `structname` and `pkgname`.  The sas_to_parquet command
generates a Go program called 'convert_data.go' that you can use to
perform the data conversion.

The parquet file will be written to the specified destination
directory, which in the above example is the current working
directory.  The parquet file name will be based on the SAS file name,
e.g. in the above example it will be 'mydata.parquet'.

To facilitate reading the Parquet file into Go using the parquet-go
package, a Go struct definition will be written to the directory
specified by 'mypackage' above.  See the `sas_to_parquet_check.go`
script to see how to read the file into Go using these struct
definitions.

## Testing

Automated testing is implemented against the Stata files used to test
the pandas Stata reader (for versions 115+):

https://github.com/pydata/pandas/tree/master/pandas/io/tests/data

A CSV data file for testing is generated by the `gendat.go` script.
There are scripts `make.sas` and `make.stata` in the test directory
that generate SAS and Stata files for testing.  SAS and Stata software
are required to run these scripts.  The generated files are provided
in the `test_files/data` directory, so `go test` can be run without
having access to SAS or Stata.

The `columnize_test.go` and `stattocsv_test.go` scripts test the
commands against stored output.

## Feedback

Please file an issue if you encounter a file that is not properly
handled.  If possible, share the file that causes the problem.
package main

// Convert a binary SAS7BDAT or Stata dta file to a CSV file.  The CSV
// contents are sent to standard output.  Date variables are returned
// as numeric values with interpretation depending on the date format
// (e.g. it may be the number of days since January 1, 1960).

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/dominodatalab/datareader"
)

func main() {

	if len(os.Args) == 1 {
		fmt.Printf("usage: %s filename\n", os.Args[0])
		return
	}

	fname := os.Args[1]
	f, err := os.Open(fname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "unable to close file %s: %v\n", fname, err)
		}
	}()

	// Determine the file type
	fl := strings.ToLower(fname)
	filetype := ""
	if strings.HasSuffix(fl, "sas7bdat") {
		filetype = "sas"
	} else if strings.HasSuffix(fl, "dta") {
		filetype = "stata"
	} else {
		fmt.Fprintf(os.Stderr, "%s file cannot be read", fname)
		return
	}

	// Get a reader for either a Stata or SAS file
	var rdr datareader.StatfileReader
	switch filetype {
	case "sas":
		sas, err := datareader.NewSAS7BDATReader(f)
		if err != nil {
			panic(err)
		}
		sas.ConvertDates = true
		sas.TrimStrings = true
		rdr = sas
	case "stata":
		stata, err := datareader.NewStataReader(f)
		if err != nil {
			panic(err)
		}
		stata.ConvertDates = true
		stata.InsertCategoryLabels = true
		stata.InsertStrls = true
		rdr = stata
	}

	err = datareader.ToCsv(rdr, 1000, csv.NewWriter(os.Stdout))
	if err != nil {
		panic(err)
	}
}

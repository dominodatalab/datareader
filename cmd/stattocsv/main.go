package main

// Convert a binary SAS7BDAT or Stata dta file to a CSV file.  The CSV
// contents are sent to standard output.  Date variables are returned
// as numeric values with interpretation depending on the date format
// (e.g. it may be the number of days since January 1, 1960).

import (
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
		os.Stderr.WriteString(fmt.Sprintf("%v\n", err))
		return
	}
	defer f.Close()

	// Determine the file type
	fl := strings.ToLower(fname)
	filetype := ""
	if strings.HasSuffix(fl, "sas7bdat") {
		filetype = "sas"
	} else if strings.HasSuffix(fl, "dta") {
		filetype = "stata"
	} else {
		os.Stderr.WriteString(fmt.Sprintf("%s file cannot be read", fname))
		return
	}

	// Get a reader for either a Stata or SAS file
	var rdr datareader.StatfileReader
	if filetype == "sas" {
		sas, err := datareader.NewSAS7BDATReader(f)
		if err != nil {
			panic(err)
		}
		sas.ConvertDates = true
		sas.TrimStrings = true
		rdr = sas
	} else if filetype == "stata" {
		stata, err := datareader.NewStataReader(f)
		if err != nil {
			panic(err)
		}
		stata.ConvertDates = true
		stata.InsertCategoryLabels = true
		stata.InsertStrls = true
		rdr = stata
	}

	datareader.DoConversion(rdr)
}

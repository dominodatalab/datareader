package datareader

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"
)

func DoConversion(rdr StatfileReader) {

	w := csv.NewWriter(os.Stdout)

	ncol := len(rdr.ColumnNames())
	if err := w.Write(rdr.ColumnNames()); err != nil {
		panic(err)
	}

	row := make([]string, ncol)

	for {
		chunk, err := rdr.Read(1000)
		if err != nil && err != io.EOF {
			panic(err)
		} else if chunk == nil || err == io.EOF {
			break
		}

		for j := 0; j < len(chunk); j++ {
			chunk[j] = chunk[j].UpcastNumeric()
		}

		nrow := chunk[0].Length()

		numbercols := make([][]float64, ncol)
		stringcols := make([][]string, ncol)
		timecols := make([][]time.Time, ncol)

		missing := make([][]bool, ncol)

		for j := 0; j < ncol; j++ {
			missing[j] = chunk[j].Missing()
			dcol := chunk[j].Data()
			switch v := dcol.(type) {
			case []time.Time:
				timecols[j] = v
			case []float64:
				numbercols[j] = v
			case []string:
				stringcols[j] = v
			default:
				panic(fmt.Sprintf("unknown type: %T", dcol))
			}
		}

		for i := range nrow {
			for j := range ncol {
				if numbercols[j] != nil {
					if missing[j] == nil || !missing[j][i] {
						row[j] = fmt.Sprintf("%f", numbercols[j][i])
					} else {
						row[j] = ""
					}
				} else if stringcols[j] != nil {
					if missing[j] == nil || !missing[j][i] {
						row[j] = stringcols[j][i]
					} else {
						row[j] = ""
					}
				} else if timecols[j] != nil {
					if missing[j] == nil || !missing[j][i] {
						row[j] = fmt.Sprintf("%v", timecols[j][i])
					} else {
						row[j] = ""
					}
				}
			}
			if err := w.Write(row); err != nil {
				panic(err)
			}
		}
	}

	w.Flush()
}

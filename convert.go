package datareader

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"time"
)

func ToCsv(rdr StatfileReader, rows int, w *csv.Writer) error {
	ncol := len(rdr.ColumnNames())
	if err := w.Write(rdr.ColumnNames()); err != nil {
		return err
	}

	row := make([]string, ncol)

	for {
		chunk, err := rdr.Read(rows)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		} else if chunk == nil || errors.Is(err, io.EOF) {
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

		for j := range ncol {
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
				return fmt.Errorf("unknown type: %T", dcol)
			}
		}

		for i := range nrow {
			for j := range ncol {
				switch {
				case numbercols[j] != nil:
					if missing[j] == nil || !missing[j][i] {
						row[j] = fmt.Sprintf("%f", numbercols[j][i])
					} else {
						row[j] = ""
					}
				case stringcols[j] != nil:
					if missing[j] == nil || !missing[j][i] {
						row[j] = stringcols[j][i]
					} else {
						row[j] = ""
					}
				case timecols[j] != nil:
					if missing[j] == nil || !missing[j][i] {
						row[j] = timecols[j][i].String()
					} else {
						row[j] = ""
					}
				}
			}
			if werr := w.Write(row); werr != nil {
				return werr
			}
		}

		// less rows came back than requested, so might be premature end to file
		if nrow < rows {
			break
		}
	}

	w.Flush()
	return nil
}

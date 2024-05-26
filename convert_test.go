package datareader_test

import (
	"bytes"
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"

	"github.com/dominodatalab/datareader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testFile struct {
	path string
}

func TestToCsvConvertsSAS(t *testing.T) {
	files, err := filepath.Glob("test_files/data/*.sas7bdat")
	require.NoError(t, err)

	testcases := map[string]testFile{}
	for _, f := range files {
		testcases["converts "+f] = testFile{path: f}
	}

	for name, tt := range testcases {
		t.Run(name, func(t *testing.T) {
			f, err := os.Open(tt.path)
			require.NoError(t, err)
			defer f.Close()

			sas, err := datareader.NewSAS7BDATReader(f)
			require.NoError(t, err)
			sas.ConvertDates = true
			sas.TrimStrings = true

			buf := new(bytes.Buffer)
			w := csv.NewWriter(buf)
			err = datareader.ToCsv(sas, 1000, w)

			require.NoError(t, err)

			r := csv.NewReader(buf)
			records, err := r.ReadAll()
			assert.NoError(t, err)
			assert.NotEmpty(t, records)
		})
	}
}

func TestToCsvConvertsStata(t *testing.T) {
	files, err := filepath.Glob("test_files/data/*.dta")
	require.NoError(t, err)

	testcases := map[string]testFile{}
	for _, f := range files {
		testcases["converts "+f] = testFile{path: f}
	}

	for name, tt := range testcases {
		t.Run(name, func(t *testing.T) {
			f, err := os.Open(tt.path)
			require.NoError(t, err)
			defer f.Close()

			stata, err := datareader.NewStataReader(f)
			require.NoError(t, err)
			stata.ConvertDates = true
			stata.InsertCategoryLabels = true
			stata.InsertStrls = true

			buf := new(bytes.Buffer)
			w := csv.NewWriter(buf)
			err = datareader.ToCsv(stata, 1000, w)
			require.NoError(t, err)

			r := csv.NewReader(buf)
			records, err := r.ReadAll()
			assert.NoError(t, err)
			assert.NotEmpty(t, records)
		})
	}
}

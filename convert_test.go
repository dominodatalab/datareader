package datareader_test

import (
	"bytes"
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"
	"time"

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

func TestToCsvConvertsTruncatedSAS(t *testing.T) {
	// project.sas7bdat is normally 1MB+, but wrote just the first 128KB to a file like:
	// { head -c 131072 >project_incomplete.sas7bdat; } < project.sas7bdat
	// also useful for generating many dummy rows
	// yes ",,0.000000,,," | head -n "40000" > "test.csv"
	f, err := os.Open("test_files/data/project_incomplete.sas7bdat")
	require.NoError(t, err)
	defer f.Close()

	sas, err := datareader.NewSAS7BDATReader(f)
	require.NoError(t, err)
	sas.ConvertDates = true
	sas.TrimStrings = true

	// verify all the file header data processed correctly
	// these values come from complete file project.sas7bdat
	assert.Equal(t,
		"PROJECT                         \x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00",
		sas.Name)
	assert.Equal(t, "DATA    ", sas.FileType)
	assert.Equal(t, "utf-8", sas.FileEncoding)
	assert.True(t, sas.U64)
	assert.Equal(t, "2.6.32-754.43.1.", sas.OSType)
	assert.Equal(t, "x86_64", sas.OSName)
	assert.Equal(t, "9.0401M6", sas.SASRelease)
	assert.Equal(t, "Linux", sas.ServerType)
	assert.Equal(t, 46641, sas.RowCount())
	assert.Equal(t, sas.ColumnLabels(), []string{
		"Type of alteration or repair",
		"Household member performed alteration or repair",
		"Cost of alteration or repair",
		"Edit flag for RAS",
		"Edit flag for RAD",
		"Control number",
	})
	assert.Equal(t, sas.ColumnNames(), []string{"RAS", "RAH", "RAD", "JRAS", "JRAD", "CONTROL"})
	assert.Equal(t, sas.ColumnTypes(), []datareader.ColumnTypeT{
		datareader.SASStringType,
		datareader.SASStringType,
		datareader.SASNumericType,
		datareader.SASStringType,
		datareader.SASStringType,
		datareader.SASStringType,
	})
	// Timestamp is epoch 01/01/1960
	tv := float64(1969085979.342952)
	ts := time.Date(1960, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(tv) * time.Second)
	assert.Equal(t, ts, sas.DateCreated)
	assert.Equal(t, ts, sas.DateModified)

	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)
	err = datareader.ToCsv(sas, 1000, w)
	require.NoError(t, err)

	r := csv.NewReader(buf)
	records, err := r.ReadAll()
	assert.NoError(t, err)
	assert.NotEmpty(t, records)

	// there are 46641 records in the file, but b/c the file is truncated there are only 2605 parsed rows
	// do a sanity check on the last correctly read row
	assert.Equal(t, []string{"62", "2", "1500.000000", "", "", "186336030133"}, records[2604])

	// NOTE: bugfix shows that only the rows that could be read are captured in the file - the write aborts once it reaches incomplete data
	assert.Equal(t, 2605, len(records))
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

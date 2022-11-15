package util

import (
	"encoding/csv"
	"net/http"
)

func OpenCsvFile(path string) ([][]string, error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	reader.LazyQuotes = true

	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	all, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return all, nil
}

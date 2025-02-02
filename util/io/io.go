package io

import (
	"os"

	"github.com/gocarina/gocsv"
)

func ReadCsvInto(path string, obj any) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := gocsv.UnmarshalFile(file, obj); err != nil {
		return err
	}
	return nil
}

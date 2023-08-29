package csv

import (
	"fmt"
	"os"
)

type stringer interface {
	String() string
}

func GenerateCSV[T stringer](data []*T) (string, error) {
	file, err := os.CreateTemp("", "data-*.csv")
	if err != nil {
		return "", err
	}
	defer file.Close()
	for _, value := range data {
		_, err := fmt.Fprintln(file, value)
		if err != nil {
			return "", err
		}
	}
	return file.Name(), nil
}

package utils

import (
	"errors"
)

func New(err string) error {
	return errors.New(err)
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

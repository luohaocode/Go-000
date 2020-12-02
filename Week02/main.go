package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func dao() (int, error) {
	return 0, errors.WithMessage(sql.ErrNoRows, "cannnot find tmp")
}

func biz() (int, error) {
	return dao()
}

func main() {
	tmp, err := biz()
	if err != nil {
		fmt.Printf("original error is %T %v\n", errors.Cause(err), errors.Cause(err))
		fmt.Printf("Stack trace: %+v\n", err)
	}
	fmt.Println(tmp)
	os.Exit(1)
}

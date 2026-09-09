package main

import (
	"fmt"

	"github.com/OmarLP/api_module/pkg/bootstrap"
)

func main() {
	db, err := bootstrap.DBConnection()
	if err != nil {
		fmt.Errorf("error connecting to database: %w", err)
	}

	fmt.Println("connection succesfully", db)
}

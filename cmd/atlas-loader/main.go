package main

import (
	"fmt"
	"io"
	"os"

	"study-golang-backend/internal/domain/entity"

	"ariga.io/atlas-provider-gorm/gormschema"
)

func main() {
	loader := gormschema.New("postgres")

	// Load entities struct
	stmts, err := loader.Load(
		&entity.User{},
		&entity.Product{},
		&entity.Tag{},
		&entity.Category{},
		&entity.CartItems{},
	)

	if err != nil {
		fmt.Println("Error loading entities:", err)
		os.Exit(1)
	}

	io.WriteString(os.Stdout, stmts)
}

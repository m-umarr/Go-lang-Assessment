package main

import (
	"github.com/organization_api/pkg"
	db "github.com/organization_api/pkg/database"
)

func main() {
	// Connect to the database.
	err := db.Connect()
	if err != nil {
		panic(err)
	}

	pkg.Init()
}

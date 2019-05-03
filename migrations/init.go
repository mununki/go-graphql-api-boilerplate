package main

import (
	"github.com/mattdamon108/go-graphql-api-boilerplate/db"
	"github.com/mattdamon108/go-graphql-api-boilerplate/model"
)

func main() {
	d, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}

	defer d.Close()

	d.DropTableIfExists(&model.User{})
	d.CreateTable(&model.User{})
}

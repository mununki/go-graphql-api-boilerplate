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

	defer d.DB.Close()

	d.DB.DropTableIfExists(&model.User{})
	d.DB.CreateTable(&model.User{})
}

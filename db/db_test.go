package db_test

import (
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/mattdamon108/go-graphql-api-boilerplate/db"
	"github.com/mattdamon108/go-graphql-api-boilerplate/model"
)

func TestDBSelect(t *testing.T) {
	db, err := db.ConnectDB()
	if err != nil {
		t.Error("Failed to connect DB")
	}

	user := model.UserAndSocial{}
	found, err := db.
		From("user").
		Join(goqu.T("user_social"), goqu.On(goqu.I("user.id").Eq(goqu.I("user_social.user_id")))).
		Where(goqu.I("user.id").Eq(1)).
		ScanStruct(&user)
	if err != nil {
		t.Error(err)
	}
	if !found {
		t.Errorf("Can't found %v", err)
	} else {
		t.Log(user)
	}
}

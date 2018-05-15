package tests

import (
	"github.com/HackIllinois/api-user/config"
	"github.com/HackIllinois/api-user/database"
	"github.com/HackIllinois/api-user/models"
	"github.com/HackIllinois/api-user/service"
	"reflect"
	"testing"
)

/*
	Initialize databse with test user info
*/
func SetupTestDB(t *testing.T) {
	err := database.Insert("info", &models.UserInfo{
		ID:       "testid",
		Username: "testusername",
		Email:    "testemail@domain.com",
	})

	if err != nil {
		t.Fatal(err)
	}
}

/*
	Drop test database
*/
func CleanupTestDB(t *testing.T) {
	session := database.GetSession()
	defer session.Close()

	err := session.DB(config.USER_DB_NAME).DropDatabase()

	if err != nil {
		t.Fatal(err)
	}
}

/*
	Service level test for getting user info from database
*/
func TestGetUserInfoService(t *testing.T) {
	SetupTestDB(t)

	user_info, err := service.GetUserInfo("testid")

	if err != nil {
		t.Fatal(err)
	}

	expected_info := &models.UserInfo{
		ID:       "testid",
		Username: "testusername",
		Email:    "testemail@domain.com",
	}

	if !reflect.DeepEqual(user_info, expected_info) {
		t.Errorf("Wrong user info. Expected %v, got %v", expected_info, user_info)
	}

	CleanupTestDB(t)
}

/*
	Service level test for setting user info in the database
*/
func TestSetUserInfoService(t *testing.T) {
	SetupTestDB(t)

	err := service.SetUserInfo("testid2", models.UserInfo{
		ID:       "testid2",
		Username: "testusername2",
		Email:    "testemail2@domain.com",
	})

	if err != nil {
		t.Fatal(err)
	}

	user_info, err := service.GetUserInfo("testid2")

	if err != nil {
		t.Fatal(err)
	}

	expected_info := &models.UserInfo{
		ID:       "testid2",
		Username: "testusername2",
		Email:    "testemail2@domain.com",
	}

	if !reflect.DeepEqual(user_info, expected_info) {
		t.Errorf("Wrong user info. Expected %v, got %v", expected_info, user_info)
	}

	CleanupTestDB(t)
}

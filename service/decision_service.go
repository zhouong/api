package service

import (
	"errors"

	"github.com/HackIllinois/api-commons/database"
	"github.com/HackIllinois/api-decision/config"
	"github.com/HackIllinois/api-decision/models"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

var db database.MongoDatabase

func init() {
	db_connection, err := database.InitMongoDatabase(config.DECISION_DB_HOST, config.DECISION_DB_NAME)

	if err != nil {
		panic(err)
	}

	db = db_connection
}

/*
	Returns the decision associated with the given user id
*/
func GetDecision(id string) (*models.DecisionHistory, error) {
	query := bson.M{
		"id": id,
	}

	var decision models.DecisionHistory
	err := db.FindOne("decision", query, &decision)

	if err != nil {
		return nil, err
	}

	return &decision, nil
}

/*
	Updates the decision associated with the given user id
	If a decision doesn't exist it will be created
*/
func UpdateDecision(id string, decision models.Decision) error {
	err := validate.Struct(decision)

	if err != nil {
		return err
	}

	if decision.Status == "ACCEPTED" && decision.Wave == 0 {
		return errors.New("Must set a wave for accepted attendee")
	} else if decision.Status != "ACCEPTED" && decision.Wave != 0 {
		return errors.New("Cannot set a wave for non-accepted attendee")
	}

	decision_history, err := GetDecision(id)

	if err != nil {
		if err == mgo.ErrNotFound {
			decision_history = &models.DecisionHistory{
				ID: id,
			}
		} else {
			return err
		}
	}

	decision_history.Finalized = decision.Finalized
	decision_history.Status = decision.Status
	decision_history.Wave = decision.Wave
	decision_history.History = append(decision_history.History, decision)
	decision_history.Reviewer = decision.Reviewer
	decision_history.Timestamp = decision.Timestamp

	selector := bson.M{
		"id": id,
	}

	err = db.Update("decision", selector, &decision_history)

	if err == mgo.ErrNotFound {
		err = db.Insert("decision", &decision_history)
	}

	return err
}

/*
	Checks if a decision with the provided id exists.
*/
func HasDecision(id string) (bool, error) {
	_, err := GetDecision(id)

	if err == nil {
		return true, nil
	} else if err == mgo.ErrNotFound {
		return false, nil
	} else {
		return false, err
	}
}
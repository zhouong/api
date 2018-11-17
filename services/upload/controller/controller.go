package controller

import (
	"encoding/json"
	"github.com/HackIllinois/api/common/errors"
	"github.com/HackIllinois/api/services/upload/service"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
)

func SetupController(route *mux.Route) {
	router := route.Subrouter()

	router.Handle("/resume/upload/", alice.New().ThenFunc(GetUpdateUserResume)).Methods("GET")
	router.Handle("/resume/", alice.New().ThenFunc(GetCurrentUserResume)).Methods("GET")
	router.Handle("/resume/{id}/", alice.New().ThenFunc(GetUserResume)).Methods("GET")
}

/*
	Endpoint to get a specified user's resume
*/
func GetUserResume(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	resume, err := service.GetUserResumeLink(id)

	if err != nil {
		panic(errors.UnprocessableError(err.Error()))
	}

	json.NewEncoder(w).Encode(resume)
}

/*
	Endpoint to get the current user's resume
*/
func GetCurrentUserResume(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("HackIllinois-Identity")

	resume, err := service.GetUserResumeLink(id)

	if err != nil {
		panic(errors.UnprocessableError(err.Error()))
	}

	json.NewEncoder(w).Encode(resume)
}

/*
	Endpoint to update the specified user's resume
*/
func GetUpdateUserResume(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("HackIllinois-Identity")

	resume, err := service.GetUpdateUserResumeLink(id)

	if err != nil {
		panic(errors.UnprocessableError(err.Error()))
	}

	json.NewEncoder(w).Encode(resume)
}

package api

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/orangejohny/SBWeb/internal/model"
)

// Config for api package. Address is domain name of server
type Config struct {
	Address      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// StartServer creates and runs API server
func StartServer(cfg Config, m *model.Model) {
	r := mux.NewRouter()
	r.Host(cfg.Address)

	r.Handle("/ads", ReadMultipleAds(m)).Methods("GET")
	r.Handle("/ads/{id:[0-9]+}", ReadOneAd(m)).Methods("GET")
	r.Handle("/users/{id:[0-9]+}", ReadUserWithID(m)).Methods("GET")

	server := http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	server.ListenAndServe()
}

// ReadMultipleAds handles */ads. It responses with list of Ads. Method is GET
func ReadMultipleAds(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		offset, _ := strconv.Atoi(r.FormValue("offset"))
		limit, err := strconv.Atoi(r.FormValue("limit"))
		if err != nil {
			limit = 15 // should be configurable
		}

		ads, err := m.GetAds(limit, offset)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't take information from database", "DatabaseError", err))
			return
		}

		adsData, err := json.Marshal(ads)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't encode JSON", "JSONerror", err))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(adsData)
	})
}

// ReadOneAd handles */ads/{id:[0-9]+} with method GET. Returns one ad with ID provided from URL
func ReadOneAd(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		idStr, _ := mux.Vars(r)["id"]
		id, _ := strconv.Atoi(idStr)
		ad, err := m.GetAd(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't take information from database", "DatabaseError", err))
			return
		}

		adData, err := json.Marshal(ad)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't encode JSON", "JSONerror", err))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(adData)
	})
}

// ReadUserWithID handles */users/{id:[0-9]+} with method GET. Returns one user struct with ID provided from URL
// TODO implement parameter show_ads
func ReadUserWithID(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		idStr, _ := mux.Vars(r)["id"]
		id, _ := strconv.Atoi(idStr)
		user, err := m.GetUserWithID(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't take information from database", "DatabaseError", err))
			return
		}

		userData, err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't encode JSON", "JSONerror", err))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(userData)
	})
}

// AddNewUser handles */users/new with method POST.
func AddNewUser(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle("Can't parse request body", "RequestFormParseError", err))
			return
		}

		var user model.User
		decoder := schema.NewDecoder()
		err = decoder.Decode(&user, r.Form)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't decode request body", "RequestFormDecodeError", err))
			return
		}

		_, err = govalidator.ValidateStruct(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle("Data didn't passed validation", "RequestDataValidError", err))
			return
		}

		if user.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle("Password wasn't received", "NoPasswordError", errors.New("Need Password to create new user")))
			return
		}

		h := sha256.New()
		h.Write([]byte(user.Password))
		user.Password = string(h.Sum(nil))

		id, err := m.NewUser(&user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't create new user", "UserCreatingError", err))
			return
		}

		user.ID = id
		user.Password = ""
		userData, err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't encode JSON", "JSONerror", err))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(userData)
	})
}

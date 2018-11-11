package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"bmstu.codes/developers34/SBWeb/internal/model"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"golang.org/x/crypto/bcrypt"
)

// StartServer creates and runs API server
func StartServer(cfg Config, m *model.Model) {
	r := mux.NewRouter()
	//r.Host(cfg.Address)

	r.Handle("/ads", readMultipleAds(m)).Methods("GET")
	r.Handle("/ads/{id:[0-9]+}", readOneAd(m)).Methods("GET")
	r.Handle("/users/{id:[0-9]+}", readUserWithID(m)).Methods("GET")

	r.Handle("/users/new", userCreatePage(m)).Methods("POST")

	r.Handle("/users/login", userLoginPage(m)).Methods("POST")
	r.Handle("/users/logout", userLogoutPage(m)).Methods("POST")

	r.Handle("/users/profile", checkCookieMiddleware(m, userProfilePage(m))).Methods("GET")
	r.Handle("/users/profile", checkCookieMiddleware(m, userUpdatePage(m))).Methods("POST")
	r.Handle("/users/profile", checkCookieMiddleware(m, userDeletePage(m))).Methods("DELETE")

	r.Handle("/ads/new", checkCookieMiddleware(m, adCreatePage(m))).Methods("POST")
	r.Handle("/ads/edit/{id:[0-9]+}", checkCookieMiddleware(m, adUpdatePage(m))).Methods("POST")
	r.Handle("/ads/delete/{id:[0-9]+}", checkCookieMiddleware(m, adDeletePage(m))).Methods("DELETE")

	RT, err1 := time.ParseDuration(cfg.ReadTimeout)
	WT, err2 := time.ParseDuration(cfg.WriteTimeout)
	IT, err3 := time.ParseDuration(cfg.IdleTimeout)
	if err1 != nil || err2 != nil || err3 != nil {
		log.Fatalln("Can't parse API config")
	}

	server := http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  RT,
		WriteTimeout: WT,
		IdleTimeout:  IT,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln(err.Error())
	}
}

// readMultipleAds handles */ads. It responses with list of Ads. Method is GET
func readMultipleAds(m *model.Model) http.Handler {
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

// readOneAd handles */ads/{id:[0-9]+} with method GET. Returns one ad with ID provided from URL
func readOneAd(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		idStr, _ := mux.Vars(r)["id"]
		id, _ := strconv.ParseInt(idStr, 10, 64)
		ad, err := m.GetAd(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't take information from database", "DatabaseError", err))
			return
		}

		if ad.Description == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle("No ad with such ID", "NoSuchAd", errors.New("No ad with such ID")))
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

// readUserWithID handles */users/{id:[0-9]+} with method GET. Returns one user struct with ID provided from URL
// if parameter show_ads == true it return list of ads of such user
func readUserWithID(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		idStr, _ := mux.Vars(r)["id"]
		id, _ := strconv.ParseInt(idStr, 10, 64)
		user, err := m.GetUserWithID(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't take information from database", "DatabaseError", err))
			return
		}

		if user.Email == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle("No user with such ID", "NoSuchUser", errors.New("No user with such ID")))
			return
		}

		userData, err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't encode JSON", "JSONerror", err))
			return
		}

		showAds := r.FormValue("show_ads")
		if showAds == "true" {
			ads, err := m.GetAdsOfUser(id)
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
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(userData)
	})
}

// userCreatePage handles */users/new with method POST.
func userCreatePage(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		// trying to parse form
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle("Can't parse request body", "RequestFormParseError", err))
			return
		}

		// get info about user from request
		var user model.User
		decoder := schema.NewDecoder()
		err = decoder.Decode(&user, r.Form)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't decode request body", "RequestFormDecodeError", err))
			return
		}

		// check data is not null explicitly
		if user.FirstName == "" || user.LastName == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle("No required information", "NoInfoError", errors.New("Need more info to create new user")))
			return
		}

		// validate incoming data; it also checks email and password are not null
		_, err = govalidator.ValidateStruct(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle("Data didn't passed validation", "RequestDataValidError", err))
			return
		}

		// make hash from incoming password
		hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		user.Password = string(hash)

		// add user to database
		id, err := m.NewUser(&user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't create new user", "UserCreatingError", err))
			return
		}

		// send user id as a response
		userData, err := json.Marshal(struct {
			ID  int64
			Ref string
		}{
			ID:  id,
			Ref: "/users/" + strconv.FormatInt(id, 10),
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't encode JSON", "JSONerror", err))
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(userData)
	})
}

// userLoginPage handles */users/login with method POST
func userLoginPage(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		// trying to parse form
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle("Can't parse request body", "RequestFormParseError", err))
			return
		}

		// get info about user from request
		var user model.User
		decoder := schema.NewDecoder()
		err = decoder.Decode(&user, r.Form)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't decode request body", "RequestFormDecodeError", err))
			return
		}

		// validate incoming data; it also checks email and password are not null
		_, err = govalidator.ValidateStruct(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle("Data didn't passed validation", "RequestDataValidError", err))
			return
		}

		// trying to find user with such email in database
		userFromDB, err := m.GetUserWithEmail(user.Email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't take information from database", "DatabaseError", err))
			return
		}

		if userFromDB.Email == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle("No user with such email", "NoSuchUser", errors.New("No user with such email")))
			return
		}

		// check if password is valid

		if err = bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(user.Password)); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle("Login or password is incorrect", "BadAuth", errors.New("Login or password is incorrect")))
			return
		}

		// create new session for user
		sess, err := m.CreateSession(&model.Session{
			ID:        userFromDB.ID,
			Login:     user.Email,
			UserAgent: r.UserAgent(),
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't create new session", "SessionCreateError", err))
			return
		}

		// set cookie with session ID
		cookie := http.Cookie{
			Name:    "session_id",
			Value:   sess.ID,
			Expires: time.Now().Add(1 * time.Hour), // should be configureable
		}

		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusFound)
	})
}

// userLogoutPage handles */users/logout with method POST
func userLogoutPage(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("session_id")
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusFound)
			return
		}

		m.DeleteSession(&model.SessionID{
			ID: session.Value,
		})

		// delete cookie
		session.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, session)

		w.WriteHeader(http.StatusFound)
	})
}

// userUpdatePage handles */users/profile with method POST
func userUpdatePage(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		// trying to parse form
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle("Can't parse request body", "RequestFormParseError", err))
			return
		}

		// get info about user from request
		var user model.User
		decoder := schema.NewDecoder()
		err = decoder.Decode(&user, r.Form)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't decode request body", "RequestFormDecodeError", err))
			return
		}

		// check data is not null explicitly
		if user.FirstName == "" || user.LastName == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle("No required information", "NoInfoError", errors.New("Need more info to update user")))
			return
		}

		user.ID = getIDfromCookie(m, r)

		_, err = m.EditUser(&user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't update user", "UserUpdatingError", err))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// userProfilePage handles */users/profile with method GET
func userProfilePage(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		user, err := m.GetUserWithID(getIDfromCookie(m, r))
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

// userDeletePage handles */users/profile with method DELETE
func userDeletePage(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		_, err := m.RemoveUser(getIDfromCookie(m, r))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't delete user from database", "DatabaseError", err))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// adCreatePage handles */ads/new with method POST
func adCreatePage(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		// trying to parse form
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle("Can't parse request body", "RequestFormParseError", err))
			return
		}

		// get info about ad from request
		var ad model.AdItem
		decoder := schema.NewDecoder()
		err = decoder.Decode(&ad, r.Form)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't decode request body", "RequestFormDecodeError", err))
			return
		}

		// validate incoming data
		_, err = govalidator.ValidateStruct(&ad)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle("Data didn't passed validation", "RequestDataValidError", err))
			return
		}

		ad.UserID = getIDfromCookie(m, r)

		// add ad to database
		id, err := m.NewAd(&ad)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't create new ad", "AdCreatingError", err))
			return
		}

		// send ad id as a response
		adData, err := json.Marshal(struct {
			ID  int64
			Ref string
		}{
			ID:  id,
			Ref: "/ads/" + strconv.FormatInt(id, 10),
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't encode JSON", "JSONerror", err))
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(adData)
	})
}

// adUpdatePage handles */ads/edit/{id:[0-9]+} with method POST
func adUpdatePage(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		idStr, _ := mux.Vars(r)["id"]
		id, _ := strconv.ParseInt(idStr, 10, 64)

		// trying to parse form
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle("Can't parse request body", "RequestFormParseError", err))
			return
		}

		// get info about ad from request
		var ad model.AdItem
		decoder := schema.NewDecoder()
		err = decoder.Decode(&ad, r.Form)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't decode request body", "RequestFormDecodeError", err))
			return
		}

		// validate incoming data
		_, err = govalidator.ValidateStruct(&ad)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle("Data didn't passed validation", "RequestDataValidError", err))
			return
		}

		ad.User.ID = getIDfromCookie(m, r)
		ad.ID = id

		adFromDatabase, err := m.GetAd(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't take information from database", "DatabaseError", err))
			return
		}

		if ad.User.ID != adFromDatabase.User.ID {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle("Trying to change ad that created by other user", "AdUpdateError", errors.New("You can change ads that created only by yourself")))
			return
		}

		_, err = m.EditAd(&ad)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't update ad", "AdUpdatingError", err))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// adDeletePage handles */ads/delete/{id:[0-9]+} with method DELETE
func adDeletePage(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		idStr, _ := mux.Vars(r)["id"]
		id, _ := strconv.ParseInt(idStr, 10, 64)

		ownerIDfromCookie := getIDfromCookie(m, r)
		adFromDatabase, err := m.GetAd(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't take information from database", "DatabaseError", err))
			return
		}

		if ownerIDfromCookie != adFromDatabase.User.ID {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle("Trying to change ad that created by other user", "AdUpdateError", errors.New("You can change ads that created only by yourself")))
			return
		}

		_, err = m.RemoveAd(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle("Can't delete ad", "AdDeleteError", err))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

// File api.go contains handlers of URL addresses and
// function that initiates API server.

package api

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"bmstu.codes/developers34/SBWeb/pkg/model"
	"golang.org/x/crypto/bcrypt"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// StartServer configures and runs API server. It's always returns channel with errors
// to monitor state of server which is running in other goroutine.
func StartServer(cfg Config, m *model.Model) (*http.Server, chan error) {
	r := mux.NewRouter()
	//r.Host(cfg.Address)

	// set handlers
	r.Handle("/ads", readMultipleAds(m)).Methods("GET")
	r.Handle("/ads/{id:[0-9]+}", readOneAd(m)).Methods("GET")

	r.Handle("/users/{id:[0-9]+}", readUserWithID(m)).Methods("GET")

	r.Handle("/users/new", userCreatePage(m)).Methods("POST")
	r.Handle("/users/login", checkConnSM(m, logRequestMiddleware(m, userLoginPage(m)))).Methods("POST")
	r.Handle("/users/logout", checkConnSM(m, userLogoutPage(m))).Methods("POST", "DELETE")

	r.Handle("/users/profile",
		checkConnSM(m, checkCookieMiddleware(m, userProfilePage(m)))).Methods("GET")
	r.Handle("/users/profile",
		checkConnSM(m, checkCookieMiddleware(m, userUpdatePage(m)))).Methods("POST")
	r.Handle("/users/profile",
		checkConnSM(m, checkCookieMiddleware(m, userDeletePage(m)))).Methods("DELETE")

	r.Handle("/ads/new",
		checkConnSM(m, checkCookieMiddleware(m, adCreatePage(m)))).Methods("POST")
	r.Handle("/ads/edit/{id:[0-9]+}",
		checkConnSM(m, checkCookieMiddleware(m, adUpdatePage(m)))).Methods("POST")
	r.Handle("/ads/delete/{id:[0-9]+}",
		checkConnSM(m, checkCookieMiddleware(m, adDeletePage(m)))).Methods("DELETE")

	r.Handle("/images/{filename}", sendImage(m)).Methods("GET")

	ch := make(chan error, 1)

	// parse config times
	RT, err1 := time.ParseDuration(cfg.ReadTimeout)
	WT, err2 := time.ParseDuration(cfg.WriteTimeout)
	IT, err3 := time.ParseDuration(cfg.IdleTimeout)
	if err1 != nil || err2 != nil || err3 != nil {
		ch <- errors.New("Can't parse API config")
		log.Println("Can't parse API config")
		return nil, ch
	}

	// create server
	server := http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  RT,
		WriteTimeout: WT,
		IdleTimeout:  IT,
	}

	// create directory for images if not exist
	_ = os.Mkdir("images", 0777)
	/* if err != nil && err != os.PathError {
		ch <- err
		log.Println(err)
		return nil, ch
	} */

	// run server
	go func() {
		if err := server.ListenAndServe(); err != nil {
			ch <- err
			log.Println(err.Error())
		}
	}()

	return &server, ch
}

// readMultipleAds handles */ads with method GET. Allowed parameters are: query, limit, offset.
// Default value for offset is 0; for limit is 15. If there are no ads, sends an empty JSON array
func readMultipleAds(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		// take params from request
		// if there are some errors, then it will be handled while validation
		var params model.SearchParams
		r.ParseForm()
		decoder := schema.NewDecoder()
		decoder.Decode(&params, r.Form)

		// check if parameters are valid
		if params.Limit <= 0 {
			params.Limit = 15
		}
		if params.Offset < 0 {
			params.Offset = 0
		}

		// TODO query should have same restrictions like title
		// check if query is valid
		/* if !govalidator.IsPrintableASCII(params.Query) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(checkReq, "QueryValidError", errors.New("Bad query"),
				"Query must be printable ASCII"))
			return
		} */

		// get list of ads from DB. If there are no ads, send an empty JSON array
		ads, err := m.GetAds(&params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle(connectProvider, getInfoDBErr, err, getInfoDBMsg))
			return
		}

		// marshall list of ads to JSON format
		adsData, _ := json.Marshal(ads)

		// send data as a response
		w.WriteHeader(http.StatusOK)
		w.Write(adsData)
	})
}

// readOneAd handles */ads/{id:[0-9]+} with method GET. Returns one ad with ID provided from URL.
func readOneAd(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		// take id from url
		idStr, _ := mux.Vars(r)["id"]
		id, _ := strconv.ParseInt(idStr, 10, 64)

		// get ad from DB
		ad, err := m.GetAd(id)

		// check if ad exists
		if ad.ID == -1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(enterExID, adIDErr,
				errors.New("Client has entered wrong ID"), badIDMsg))
			return
		}

		// process error from DB
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle(connectProvider, getInfoDBErr, err, getInfoDBMsg))
			return
		}

		// marshall data to JSON format
		adData, _ := json.Marshal(ad)

		// send data as a response
		w.WriteHeader(http.StatusOK)
		w.Write(adData)
	})
}

// readUserWithID handles */users/{id:[0-9]+} with method GET. Returns one user struct with ID provided from URL.
// if parameter show_ads == true function will return list of ads of such user.
// if such user has no ads then empty JSON array will be returned.
func readUserWithID(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		// take id from url
		idStr, _ := mux.Vars(r)["id"]
		id, _ := strconv.ParseInt(idStr, 10, 64)

		// parse parameter
		showAds := r.FormValue("show_ads")
		if showAds == "true" {
			// get ads of user from DB. If user has no ads, returns empty JSON array
			ads, err := m.GetAdsOfUser(id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(apiErrorHandle(connectProvider, getInfoDBErr, err,
					getInfoDBMsg))
				return
			}

			// marshall data to JSON format
			adsData, _ := json.Marshal(ads)

			// send ads as a response
			w.WriteHeader(http.StatusOK)
			w.Write(adsData)
			return
		}

		// get user from DB
		user, err := m.GetUserWithID(id)

		// check if user exists
		if user.ID == -1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(enterExID, userIDErr,
				errors.New("Client entered wrong ID"), badIDMsg))
			return
		}

		// process error from DB
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle(connectProvider, getInfoDBErr, err, getInfoDBMsg))
			return
		}

		// marshall data to JSON format
		userData, _ := json.Marshal(user)

		// send user as a response
		w.WriteHeader(http.StatusOK)
		w.Write(userData)
	})
}

// userCreatePage handles */users/new with method POST.
// Function process incoming parameters to create new user. On success
// it will return JSON object with ID and reference to created user.
func userCreatePage(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		var isMultipartForm bool
		if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			isMultipartForm = true
		}

		// trying to parse form
		var err error
		if isMultipartForm {
			err = r.ParseMultipartForm(10 * 1024 * 1024)
		} else {
			err = r.ParseForm()
		}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(checkReq, parseFormErr, err, parseFormMsg))
			return
		}

		// get info about user from request
		var user model.User
		decoder := schema.NewDecoder()
		if isMultipartForm {
			err = decoder.Decode(&user, r.MultipartForm.Value)

		} else {
			err = decoder.Decode(&user, r.Form)
		}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(checkReq, decodeFormErr, err,
				decodeFormMsg))
			return
		}

		// check data is not null explicitly
		if user.FirstName == "" || user.LastName == "" ||
			user.Password == "" || user.Email == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(
				enterRequiredInfo,
				requiredinfoErr,
				errors.New("Client didn't sent required info"),
				requiredinfoMsg))
			return
		}

		// validate incoming data
		_, err = govalidator.ValidateStruct(&user)
		if err != nil || !govalidator.IsNumeric(user.TelNumber.String) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(validRequired, reqValidErr, err, reqValidMsg)) // err will be nil if TelNumber or About didn't passed validation
			return
		}

		// make hash from incoming password
		hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password),
			bcrypt.DefaultCost)
		user.Password = string(hash)

		// to be sure that client didn't provide avatar_address
		user.AvatarAddress.String = ""
		user.AvatarAddress.Valid = false

		// load images from request if it is possible
		if isMultipartForm {
			filenames, err := loadImages(r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(apiErrorHandle(checkImage, imgCreErr, err,
					imgCreMsg))
				return
			}
			if len(filenames) != 0 {
				user.AvatarAddress.SetValid(filenames[0])
			}
		}

		// add user to database
		id, err := m.NewUser(&user)

		// check if user exists
		if id == -1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(
				enterExEmail,
				userExErr,
				errors.New("Client tried to create user with an existing email"),
				userExMsg))
			return
		}

		// process error from DB
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle(connectProvider, addUserDBError, err,
				addUserDBMsg))
			return
		}

		// marshall data to JSON format
		userData, _ := json.Marshal(struct {
			ID  int64
			Ref string
		}{
			ID:  id,
			Ref: "/users/" + strconv.FormatInt(id, 10),
		})

		// send response
		w.WriteHeader(http.StatusCreated)
		w.Write(userData)
	})
}

// userLoginPage handles */users/login with method POST. It process incoming
// password and email to authentificate clent. On succeed it will create session,
// set cookie to response and return first name, last name, id if user agent is "Android_app".
func userLoginPage(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		// trying to parse form
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(checkReq, parseFormErr, err, parseFormMsg))
			return
		}

		// get info about user from request
		var user model.User
		decoder := schema.NewDecoder()
		err = decoder.Decode(&user, r.Form)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(checkReq, decodeFormErr, err,
				decodeFormMsg))
			return
		}

		// check data is not null explicitly
		if user.Password == "" || user.Email == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(
				enterRequiredInfoLogin,
				requiredinfoErr,
				errors.New("Client didn't sent required info"),
				requiredinfoMsgLogin))
			return
		}

		// validate incoming data
		_, err = govalidator.ValidateStruct(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(validRequired, reqValidErr, err, reqValidMsg))
			return
		}

		// trying to find user with such email in database
		userFromDB, err := m.GetUserWithEmail(user.Email)

		// check if user exists
		// empty := model.User{}
		if userFromDB.ID == -1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(enterValidAuth, badAuthErr,
				errors.New("Client has entered non-existing email"), badAuthMsg))
			return
		}

		// process error from DB
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle(connectProvider, getInfoDBErr, err, getInfoDBMsg))
			return
		}

		// check if password is valid
		if err = bcrypt.CompareHashAndPassword([]byte(userFromDB.Password),
			[]byte(user.Password)); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(enterValidAuth, badAuthErr,
				errors.New("Client has entered wrong password"), badAuthMsg))
			return
		}

		// android app don't need to set expiration
		isExpires := true
		if r.UserAgent() == "Android_app" {
			isExpires = false
		}

		// create new session for user
		sess, err := m.CreateSession(&model.Session{
			ID:        userFromDB.ID,
			Login:     user.Email,
			UserAgent: r.UserAgent(),
		}, isExpires)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle(connectProvider, sessCreErr, err, sessCreMsg))
			return
		}

		// set cookie for web-browser
		if isExpires {
			cookie := http.Cookie{
				Name:     "session_id",
				Value:    sess.ID,
				Expires:  time.Now().Add(24 * time.Hour), // TODO: should be configurable
				HttpOnly: true,
			}
			http.SetCookie(w, &cookie)

		} else {
			// set cookie  without expiration time
			cookie := http.Cookie{
				Name:     "session_id",
				Value:    sess.ID,
				HttpOnly: true,
			}
			http.SetCookie(w, &cookie)

			// send needed information to android app in JSON format
			appData, _ := json.Marshal(struct {
				// Name      string
				// Value     string `json:"session_id,"`
				ID        int64  `json:"id,"`
				FirstName string `json:"first_name,"`
				LastName  string `json:"last_name,"`
			}{
				// Name:      "session_id",
				// Value:     sess.ID,
				ID:        userFromDB.ID,
				FirstName: userFromDB.FirstName,
				LastName:  userFromDB.LastName,
			})

			w.Write(appData)
		}
	})
}

// userLogoutPage handles */users/logout with method POST. Middleware that checks
// cookie is required. Deletes current session and return status OK.
func userLogoutPage(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("session_id")
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusOK)
			return
		}

		// TODO: should handle error
		err = m.DeleteSession(&model.SessionID{
			ID: session.Value,
		})
		if err != nil {
			log.Println(err.Error())
		}

		// delete cookie
		session.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, session)

		w.WriteHeader(http.StatusOK)
	})
}

// userUpdatePage handles */users/profile with method POST. Middleware that checks
// cookie is required. Function process incoming parameters to update existing user.
// On succeed it returns status OK.
func userUpdatePage(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		var isMultipartForm bool
		if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			isMultipartForm = true
		}

		// trying to parse form
		var err error
		if isMultipartForm {
			err = r.ParseMultipartForm(10 * 1024 * 1024)
		} else {
			err = r.ParseForm()
		}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(checkReq, parseFormErr, err, parseFormMsg))
			return
		}

		// get info about user from request
		var user model.User
		decoder := schema.NewDecoder()
		if isMultipartForm {
			err = decoder.Decode(&user, r.MultipartForm.Value)

		} else {
			err = decoder.Decode(&user, r.Form)
		}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(checkReq, decodeFormErr, err,
				decodeFormMsg))
			return
		}

		// check data is not null explicitly
		if user.FirstName == "" || user.LastName == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(
				enterRequiredInfoUpdate,
				requiredinfoErr,
				errors.New("Client didn't sent required info for updating user"),
				badUpdateMsg))
			return
		}

		// validate incoming data
		_, err = govalidator.ValidateStruct(&user)
		if err != nil || !govalidator.IsNumeric(user.TelNumber.String) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(validRequiredUpdate, reqValidErr, err,
				reqValidMsg)) // err will be nil if TelNumber or About didn't passed validation
			return
		}

		// get id from request's cookie
		user.ID = getIDfromCookie(m, r)

		// check if image address not null and exists
		if user.AvatarAddress.String != "" {
			_, err = os.Stat("." + rDomain(user.AvatarAddress.String))
			if os.IsNotExist(err) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(apiErrorHandle(checkImage, imgExErr, err,
					imgExMsg))
				return
			}
		} else {
			// get user information from DB
			userFromDB, err := m.GetUserWithID(user.ID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(apiErrorHandle(connectProvider, getInfoDBErr, err, getInfoDBMsg))
				return
			}
			// delete existing avatar image
			deleteImages([]string{userFromDB.AvatarAddress.String})
		}

		// load images from request if image address is null and
		// content-type is multipart/form-data
		if isMultipartForm {
			filenames, err := loadImages(r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(apiErrorHandle(checkImage, imgCreErr, err,
					imgCreMsg))
				return
			}
			if len(filenames) != 0 {
				user.AvatarAddress.SetValid(filenames[0])
			}
		}

		// update user in DB
		_, err = m.EditUser(&user)

		// process error from DB
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle(connectProvider, updateUserDBErr, err,
				updateUserDBMsg))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// userProfilePage handles */users/profile with method GET. Requires checkCookieMiddleware.
// Returns information about logged user.
func userProfilePage(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		// get user from DB
		user, err := m.GetUserWithID(getIDfromCookie(m, r))

		// process error from DB
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle(connectProvider, getInfoDBErr, err, getInfoDBMsg))
			return
		}

		// marshall data to JSON format
		userData, _ := json.Marshal(user)

		// send response
		w.WriteHeader(http.StatusOK)
		w.Write(userData)
	})
}

// userDeletePage handles */users/profile with method DELETE. Requires checkCookieMiddleware.
// Delete current logged user. Returns status OK on succeed.
func userDeletePage(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		id := getIDfromCookie(m, r)

		// remove avatar of user
		userFromDB, err := m.GetUserWithID(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle(connectProvider, getInfoDBErr, err, getInfoDBMsg))
			return
		}
		deleteImages([]string{userFromDB.AvatarAddress.String})

		// TODO removing images of ads which were created by user

		// remove user from DB
		_, err = m.RemoveUser(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle(connectProvider, removeUserDBErr, err,
				removeUserDBMsg))
			return
		}

		// delete session
		session, _ := r.Cookie("session_id")
		m.DeleteSession(&model.SessionID{
			ID: session.Value,
		})

		// delete cookie
		session.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, session)

		w.WriteHeader(http.StatusOK)
	})
}

// adCreatePage handles */ads/new with method POST. Requires checkCookieMiddleware.
// Process parameters from request in order to create new ad. On succeed returns
// JSON object with id and reference to new ad.
func adCreatePage(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		var isMultipartForm bool
		if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			isMultipartForm = true
		}

		// trying to parse form
		var err error
		if isMultipartForm {
			err = r.ParseMultipartForm(10 * 1024 * 1024)
		} else {
			err = r.ParseForm()
		}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(checkReq, parseFormErr, err, parseFormMsg))
			return
		}

		// get info about ad from request
		var ad model.AdItem
		decoder := schema.NewDecoder()
		if isMultipartForm {
			err = decoder.Decode(&ad, r.MultipartForm.Value)

		} else {
			err = decoder.Decode(&ad, r.Form)
		}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(checkReq, decodeFormErr, err,
				decodeFormMsg))
			return
		}

		// check data is not null explicitly
		if ad.Title == "" || ad.Description == "" ||
			ad.City == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(
				enterRequiredInfoAd,
				requiredinfoErr,
				errors.New("Client didn't sent required info for ad creation"),
				requiredinfoAdMsg))
			return
		}

		// TODO custom validators for UTF-8 with some usual characters
		// validate incoming data
		/* _, err = govalidator.ValidateStruct(&ad)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(validRequiredCreateAd, reqValidErr, err,
				reqValidMsg))
			return
		} */

		// prevent client from passing this parameter
		ad.AdImages = nil
		// load images from request if it is possible
		if isMultipartForm {
			filenames, err := loadImages(r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(apiErrorHandle(checkImage, imgCreErr, err,
					imgCreMsg))
				return
			}
			if len(filenames) != 0 {
				ad.AdImages = filenames
			}
		}

		// set id from cookie
		ad.UserID = getIDfromCookie(m, r)

		// add ad to database
		// TODO: should check if ad already exists
		id, err := m.NewAd(&ad)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle(connectProvider, addAdDBErr, err,
				addAdDBMsg))
			return
		}

		// marshall data to JSON format
		adData, _ := json.Marshal(struct {
			ID  int64
			Ref string
		}{
			ID:  id,
			Ref: "/ads/" + strconv.FormatInt(id, 10),
		})

		// send response
		w.WriteHeader(http.StatusCreated)
		w.Write(adData)
	})
}

// adUpdatePage handles */ads/edit/{id:[0-9]+} with method POST. Requires checkCookieMiddleware.
// Process parameters from request to update existing user. On succeed returns status OK.
func adUpdatePage(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		// get id from url
		idStr, _ := mux.Vars(r)["id"]
		id, _ := strconv.ParseInt(idStr, 10, 64)

		var isMultipartForm bool
		if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			isMultipartForm = true
		}

		// trying to parse form
		var err error
		if isMultipartForm {
			err = r.ParseMultipartForm(10 * 1024 * 1024)
		} else {
			err = r.ParseForm()
		}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(checkReq, parseFormErr, err, parseFormMsg))
			return
		}

		// get info about ad from request
		var ad model.AdItem
		decoder := schema.NewDecoder()
		if isMultipartForm {
			err = decoder.Decode(&ad, r.MultipartForm.Value)

		} else {
			err = decoder.Decode(&ad, r.Form)
		}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(checkReq, decodeFormErr, err,
				decodeFormMsg))
			return
		}

		// check data is not null explicitly
		if ad.Title == "" || ad.Description == "" ||
			ad.City == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(
				enterRequiredInfoAd,
				requiredinfoErr,
				errors.New("Client didn't sent required info for ad update"),
				requiredInfoAdUpdateMsg))
			return
		}

		// TODO custom validators for UTF-8 with some usual characters
		// validate incoming data
		/* _, err = govalidator.ValidateStruct(&ad)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(validRequiredCreateAd, reqValidErr, err,
				reqValidMsg))
			return
		} */

		// set id to prevent updating add of other user
		ad.User.ID = getIDfromCookie(m, r)
		ad.ID = id

		// get ad from DB
		adFromDatabase, err := m.GetAd(id)

		// check if ad exists
		if adFromDatabase.ID == -1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(enterExID, adIDErr,
				errors.New("Client entered wrong ID of ad"), badIDMsg))
			return
		}

		// process DB error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle(connectProvider, getInfoDBErr, err, getInfoDBMsg))
			return
		}

		// check if client changing his ad
		if ad.User.ID != adFromDatabase.User.ID {
			w.WriteHeader(http.StatusForbidden)
			w.Write(apiErrorHandle(onlyYourAd, updateAdDBErr,
				errors.New("Client tried to change ad of other user"), onlyYourAdMsg))
			return
		}

		// check if images are not null and exist
		if ad.AdImages != nil {
			for _, image := range ad.AdImages {
				_, err = os.Stat("." + rDomain(image))
				if os.IsNotExist(err) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write(apiErrorHandle(checkImage, imgExErr, err,
						imgExMsg))
					return
				}
			}
		} else {
			// TODO delete particular images of ad
			deleteImages(adFromDatabase.AdImages)
		}

		// load images from request if existing images array is null and
		// content-type is multipart/form-data
		if isMultipartForm {
			filenames, err := loadImages(r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(apiErrorHandle(checkImage, imgCreErr, err,
					imgCreMsg))
				return
			}
			if len(filenames) != 0 {
				for _, img := range filenames {
					ad.AdImages = append(ad.AdImages, img)
				}
			}
		}

		// update ad
		_, err = m.EditAd(&ad)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle(connectProvider, updateAdDBErr, err,
				updateAdDBMsg))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// adDeletePage handles */ads/delete/{id:[0-9]+} with method DELETE.
// Requires checkCookieMiddleware. Deletes ad of current logged user.
// On succeed returns status OK.
func adDeletePage(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		// get id from url
		idStr, _ := mux.Vars(r)["id"]
		id, _ := strconv.ParseInt(idStr, 10, 64)

		// get owner from cookie and ad with such id
		ownerIDfromCookie := getIDfromCookie(m, r)
		adFromDatabase, err := m.GetAd(id)

		// check if ad exists
		if adFromDatabase.ID == -1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(enterExID, adIDErr,
				errors.New("Client entered wrong ID of ad"), badIDMsg))
			return
		}

		// process error from DB
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle(connectProvider, getInfoDBErr, err, getInfoDBMsg))
			return
		}

		// check if client is deleting exactly his ad
		if ownerIDfromCookie != adFromDatabase.User.ID {
			w.WriteHeader(http.StatusForbidden)
			w.Write(apiErrorHandle(onlyYourAd, updateAdDBErr,
				errors.New("Client tried to delete ad of other user"), onlyYourAdMsg))
			return
		}

		// delete images of such ad
		deleteImages(adFromDatabase.AdImages)

		// remove ad from DB
		_, err = m.RemoveAd(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle(connectProvider, removeAdDBErr, err,
				removeAdDBMsg))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// sendImage handles */images/{filename} with method GET.
// It returns image with such filename if it's exists.
func sendImage(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filename, _ := mux.Vars(r)["filename"]
		_, err := os.Stat("./images/" + filename)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(checkImage, noImgErr, err,
				noImgMsg))
			return
		}

		http.ServeFile(w, r, "./images/"+filename)
	})
}

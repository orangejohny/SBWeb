// Copyright Dmitry Kargashin <dkargashin3@gmail.com>
// License can be found in LICENSE file.
package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"bmstu.codes/developers34/SBWeb/internal/model"
	"golang.org/x/crypto/bcrypt"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

// StartServer creates and runs API server
func StartServer(cfg Config, m *model.Model) (*http.Server, chan error) {
	r := mux.NewRouter()
	//r.Host(cfg.Address)

	r.Handle("/ads", readMultipleAds(m)).Methods("GET")
	r.Handle("/ads/{id:[0-9]+}", readOneAd(m)).Methods("GET")

	r.Handle("/users/{id:[0-9]+}", readUserWithID(m)).Methods("GET")

	r.Handle("/users/new", userCreatePage(m)).Methods("POST")
	r.Handle("/users/login", logRequestMiddleware(m, userLoginPage(m))).Methods("POST")
	r.Handle("/users/logout", userLogoutPage(m)).Methods("POST", "DELETE")

	r.Handle("/users/profile",
		checkCookieMiddleware(m, userProfilePage(m))).Methods("GET")
	r.Handle("/users/profile",
		checkCookieMiddleware(m, userUpdatePage(m))).Methods("POST")
	r.Handle("/users/profile",
		checkCookieMiddleware(m, userDeletePage(m))).Methods("DELETE")

	r.Handle("/ads/new",
		checkCookieMiddleware(m, adCreatePage(m))).Methods("POST")
	r.Handle("/ads/edit/{id:[0-9]+}",
		checkCookieMiddleware(m, adUpdatePage(m))).Methods("POST")
	r.Handle("/ads/delete/{id:[0-9]+}",
		checkCookieMiddleware(m, adDeletePage(m))).Methods("DELETE")

	ch := make(chan error, 1)

	RT, err1 := time.ParseDuration(cfg.ReadTimeout)
	WT, err2 := time.ParseDuration(cfg.WriteTimeout)
	IT, err3 := time.ParseDuration(cfg.IdleTimeout)
	if err1 != nil || err2 != nil || err3 != nil {
		ch <- errors.New("Can't parse API config")
		log.Println("Can't parse API config")
	}

	server := http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  RT,
		WriteTimeout: WT,
		IdleTimeout:  IT,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			ch <- err
			log.Println(err.Error())
		}
	}()

	return &server, ch
}

// readMultipleAds handles */ads. It responses with list of Ads. Method is GET.
// if there are no ads, sends an empty JSON array
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

		// check if query is valid
		if !govalidator.IsPrintableASCII(params.Query) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(checkReq, decodeFormErr, errors.New("Bad query"),
				decodeFormMsg))
			return
		}

		// get list of ads from DB. If there are no ads, send an empty JSON array
		ads, err := m.GetAds(&params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle(connectProvider, getInfoDBErr, err, getInfoDBMsg))
			return
		}

		// marshall list of ads to JSON format
		adsData, err := json.Marshal(ads)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle(connectProvider, respCreErr, err, respCreMsg))
			return
		}

		// send data as a response
		w.WriteHeader(http.StatusOK)
		w.Write(adsData)
	})
}

// readOneAd handles */ads/{id:[0-9]+} with method GET. Returns one ad with ID provided from URL
func readOneAd(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		// take id from url
		idStr, _ := mux.Vars(r)["id"]
		id, _ := strconv.ParseInt(idStr, 10, 64)

		// get ad from DB
		ad, err := m.GetAd(id)

		// check if ad exists
		// empty := model.AdItem{}
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
		adData, err := json.Marshal(ad)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle(connectProvider, respCreErr, err, respCreMsg))
			return
		}

		// send data as a response
		w.WriteHeader(http.StatusOK)
		w.Write(adData)
	})
}

// readUserWithID handles */users/{id:[0-9]+} with method GET. Returns one user struct with ID provided from URL
// if parameter show_ads == true it return list of ads of such user
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
			adsData, err := json.Marshal(ads)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(apiErrorHandle(connectProvider, respCreErr, err, respCreMsg))
				return
			}

			// send ads as a response
			w.WriteHeader(http.StatusOK)
			w.Write(adsData)
			return
		}

		// get user from DB
		user, err := m.GetUserWithID(id)

		// check if user exists
		// empty := model.User{}
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
		userData, err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle(connectProvider, respCreErr, err, respCreMsg))
			return
		}

		// send user as a response
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
		if err != nil || !govalidator.IsNumeric(user.TelNumber.String) ||
			!govalidator.IsASCII(user.About.String) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(validRequired, reqValidErr, err, reqValidMsg)) // err will be nil if TelNumber or About didn't passed validation
			return
		}

		// make hash from incoming password
		hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password),
			bcrypt.DefaultCost)
		user.Password = string(hash)

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
		userData, err := json.Marshal(struct {
			ID  int64
			Ref string
		}{
			ID:  id,
			Ref: "/users/" + strconv.FormatInt(id, 10),
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle(connectProvider, respCreErr, err, respCreMsg))
			return
		}

		// send response
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
				Expires:  time.Now().Add(24 * time.Hour), // TODO: should be configureable
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
			appData, err := json.Marshal(struct {
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
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(apiErrorHandle(connectProvider, respCreErr, err, respCreMsg))
				return
			}

			w.Write(appData)
		}
	})
}

// userLogoutPage handles */users/logout with method POST
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

// userUpdatePage handles */users/profile with method POST
func userUpdatePage(m *model.Model) http.Handler {
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
			w.Write(apiErrorHandle(checkReq, decodeFormErr, err, decodeFormMsg))
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
		if err != nil || !govalidator.IsNumeric(user.TelNumber.String) ||
			!govalidator.IsASCII(user.About.String) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(validRequiredUpdate, reqValidErr, err,
				reqValidMsg)) // err will be nil if TelNumber or About didn't passed validation
			return
		}

		// get id from request's cookie
		user.ID = getIDfromCookie(m, r)

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

// userProfilePage handles */users/profile with method GET
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
		userData, err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle(connectProvider, respCreErr, err, respCreMsg))
			return
		}

		// send response
		w.WriteHeader(http.StatusOK)
		w.Write(userData)
	})
}

// userDeletePage handles */users/profile with method DELETE
func userDeletePage(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		// remove user from DB
		_, err := m.RemoveUser(getIDfromCookie(m, r))
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

// adCreatePage handles */ads/new with method POST
func adCreatePage(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		// trying to parse form
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(checkReq, parseFormErr, err, parseFormMsg))
			return
		}

		// get info about ad from request
		var ad model.AdItem
		decoder := schema.NewDecoder()
		err = decoder.Decode(&ad, r.Form)
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

		// validate incoming data
		_, err = govalidator.ValidateStruct(&ad)
		if err != nil || !govalidator.IsUTFLetter(ad.Country.String) ||
			!govalidator.IsUTFLetter(ad.SubwayStation.String) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(validRequiredCreateAd, reqValidErr, err,
				reqValidMsg))
			return
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
		adData, err := json.Marshal(struct {
			ID  int64
			Ref string
		}{
			ID:  id,
			Ref: "/ads/" + strconv.FormatInt(id, 10),
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle(connectProvider, respCreErr, err, respCreMsg))
			return
		}

		// send response
		w.WriteHeader(http.StatusCreated)
		w.Write(adData)
	})
}

// adUpdatePage handles */ads/edit/{id:[0-9]+} with method POST
func adUpdatePage(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		// get id from url
		idStr, _ := mux.Vars(r)["id"]
		id, _ := strconv.ParseInt(idStr, 10, 64)

		// trying to parse form
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(checkReq, parseFormErr, err, parseFormMsg))
			return
		}

		// get info about ad from request
		var ad model.AdItem
		decoder := schema.NewDecoder()
		err = decoder.Decode(&ad, r.Form)
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

		// validate incoming data
		_, err = govalidator.ValidateStruct(&ad)
		if err != nil || !govalidator.IsUTFLetter(ad.Country.String) ||
			!govalidator.IsUTFLetter(ad.SubwayStation.String) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiErrorHandle(validRequiredCreateAd, reqValidErr, err,
				reqValidMsg))
			return
		}

		// set id to prevent updating add of other user
		ad.User.ID = getIDfromCookie(m, r)
		ad.ID = id

		// get ad from DB
		adFromDatabase, err := m.GetAd(id)

		// check if ad exists
		// empty := model.AdItem{}
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

		// log.Printf("ad from DB: %v\nad from req:%v", adFromDatabase, ad)

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

// adDeletePage handles */ads/delete/{id:[0-9]+} with method DELETE
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
		// empty := model.AdItem{}
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

package api

import (
	"encoding/json"
	"log"
)

const (
	connectProvider         = "Connect with your API provider"
	getInfoDBMsg            = "Can't get information"
	getInfoDBErr            = "GetInfoDBError"
	respCreErr              = "ResponseCreatingError"
	respCreMsg              = "Can't encode response"
	enterExID               = "Must enter an exisiting ID"
	adIDErr                 = "NoAdWithSuchIDError"
	userIDErr               = "NoUserWithSuchID"
	userEmailErr            = "NoUserWithSuchEmail"
	badIDMsg                = "You have entered wrong ID"
	checkReq                = "Check your request"
	parseFormErr            = "RequestFormParseError"
	parseFormMsg            = "Can't parse parameters of request"
	decodeFormErr           = "RequestFormDecodeError"
	decodeFormMsg           = "Can't umarshall parameters"
	validRequired           = "First name, last name must be utf letter; password, about must be ASCII; telephone number - digits 1-9"
	reqValidErr             = "RequestDataValidError"
	reqValidMsg             = "Data didn't passed the validation"
	enterExEmail            = "Must use the unique email"
	userExErr               = "UserIsExistsError"
	userExMsg               = "User with such email is already exists"
	enterRequiredInfo       = "Enter required information (first name, last name, email, password)"
	requiredinfoErr         = "NoRequiredInfoError"
	requiredinfoMsg         = "Need more information to create new user"
	addUserDBError          = "AddUserDBError"
	addUserDBMsg            = "Can't create user"
	enterRequiredInfoLogin  = "Enter required information (email, password)"
	requiredinfoMsgLogin    = "Need more information to login"
	enterValidAuth          = "Must enter existing email/password pair"
	badAuthErr              = "BadAuth"
	badAuthMsg              = "Invalid email or password"
	sessCreErr              = "SessionCreateError"
	sessCreMsg              = "Can't create session"
	enterRequiredInfoUpdate = "Enter required information (first name, last name)"
	badUpdateMsg            = "Need more information to update user"
	validRequiredUpdate     = "First name, last name must be utf letter, about must be ASCII; telephone number - digits 1-9"
	updateUserDBErr         = "UpdateUserDBError"
	updateUserDBMsg         = "Can't update user"
	removeUserDBErr         = "RemoveUserError"
	removeUserDBMsg         = "Can't remove user"
	validRequiredCreateAd   = "Title must be utf letter or num; Country, City, Subway station must be utf letter; Description must be ASCII"
	enterRequiredInfoAd     = "Enter required information (title, city, descriprion)"
	requiredinfoAdMsg       = "Need more information to create ad"
	addAdDBErr              = "CreateAdError"
	addAdDBMsg              = "Can't create new ad"
	requiredInfoAdUpdateMsg = "Need more information to update ad"
	updateAdDBErr           = "UpdateAdError"
	updateAdDBMsg           = "Can't update ad"
	onlyYourAd              = "You can change or delete ads that created only by yourself"
	onlyYourAdMsg           = "Trying to change or delete ad that created by other user"
	removeAdDBErr           = "RemoveAdError"
	removeAdDBMsg           = "Can't remove ad"
	requiredCookie          = "You have to be authentificated to access this address"
	noCookieError           = "NoCookieError"
	noCookieMsg             = "Request doesn't have cookie"
	badCookie               = "You must have valid cookie to access this address"
	badCookieErr            = "BadCookieError"
	badCookieMsg            = "Cookie didn't passed validation"
	checkImage              = "Check your image"
	imgCreErr               = "ImageCreateError"
	imgCreMsg               = "Image must be PNG or JPEG format"
	noImgErr                = "No such image"
	noImgMsg                = "There is no image with such name"
)

type apiError struct {
	Description string `json:"description"`
	Message     string `json:"message"`
	ErrorCode   string `json:"error"`
}

func (e *apiError) Error() string {
	return e.ErrorCode
}

// apiErrorHandle creates apiError and returns JSON serialized object
func apiErrorHandle(message, errorCode string, err error, errorForClient string) []byte {
	dbErr := apiError{
		Description: errorForClient,
		Message:     message,
		ErrorCode:   errorCode,
	}
	dbErrData, _ := json.Marshal(dbErr) // should check error
	log.Println(err, string(dbErrData))
	return dbErrData
}

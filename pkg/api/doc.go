// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

/*
Package api implements handling of different URL addresses of application.
It uses data types provided by model package.

Data types

Package api uses not only default HTTP status codes but own error type
that describes happened error.

Utility types

API Error:
	description      short description of error
	message          what client should do for error resolving
	error code       unique error code

Create confirm object:
	id               identificator of created user/ad
	ref              reference to created user/ad (without base part; i.e. "/users/115")

User login confirm object for android application:
	id               identificator of logged user
	first_name       first name of user
	last_name        last name of user

User

Names of fields of JSON object which will be returned:
	id               <int64>
	first_name       <string>
	last_name        <string>
	email            <string>
	tel_number       <string>
	about            <string>
	reg_time         <string>
	avatar_address   <string>
HTTP parameters which are used to define user:
	id
	first_name
	last_name
	email              [email]
	password           [printable ascii]
	tel_number         [digits 1-9]
	about
	avatatar_address   [existing address]

Ad

Names of fields of JSON object which will be returned:
	id                 <int64>
	title              <string>
	price              <int>
	country            <string>
	city               <string>
	subway_station     <string>
	ad_images          <string array>
	owner_ad           <JSON object of user>
	description_ad     <string>
	creation_time      <string>
HTTP parameters which are used to define ad:
	id
	title
	price               [positive number]
	country
	city
	subway_station
	ad_images           [existing images addresses]
	description_ad

Interface

"base" is domain part of address (i.e. http://example.com).

Read and search multiple ads

"base/ads" address:
	method                 GET
	allowed parameters:
		query                                   search query; return only ads which contatins query in title of ad
		limit                [positive number]  maximum number of ads which will be returned
		offset               [positive number]  number of the first ad that will be returned
	return result:
		status 200           JSON array of ads
		status 400           <QueryValidError>  JSON object of API error
		status 500:
			1.           <GetInfoDBError>         JSON object of API error
			2.           <ResponseCreatingError>  JSON object of API error
If limit and/or offset aren't provided, their default values are 15 and 0.
If there is no ads then it will return empty JSON array.

Get information about particular ad

"base/ads/{id}" address:
	method                 GET
	id                     must be a digit number
	return result:
		status 200           JSON object of ad with such id
		status 400           <NoAdWithSuchIDError>   JSON object of API error
		status 500:
			1.           <GetInfoDBError>         JSON object of API error
			2.           <ResponseCreatingError>  JSON object of API error

Get information about particular user

"base/users/{id}" address:
	method                 GET
	id                     must be a digit number
	allowed parameters:
		show_ads             [true|false] if "true" then return ads of user with wuch id
	return result:
		status 200:
			1.           JSON object of user if "show_ads" isn't "true"
			2.           JSON array of ads if "show_ads" is "true"
		status 400           <NoUserWithSuchID> JSON object of API error
		status 500:
			1.           <GetInfoDBError>         JSON object of API error
			2.           <ResponseCreatingError>  JSON object of API error
If there is no ads then it will return empty JSON array.

Create new user

"base/users/new" address:
	method                 POST
	required parameters:
		first_name           [UTF letters]      first name of user
		last_name            [UTF letters]      last name of user
		email                [email]            unique email of user
		password             [printable ASCII]  password that will be used for authorization
	allowed parameters:
		tel_number           [digits 1-9]       telephone number of user
		about                                   some additional information about user
		images               [.JPEG or .png]    avatar image of user (if provided then all parameters must be in "multipart/form-data")
	return result:
		status 201           JSON object of user create confirm
		status 400:
			1.           <RequestFormParseError>  JSON object of API error
			2.           <RequestFormDecodeError> JSON object of API error
			3.           <NoRequiredInfoError>    JSON object of API error
			4.           <RequestDataValidError>  JSON object of API error
			5.           <UserIsExistsError>      JSON object of API error
		status 500:
			1.           <ImageCreateError>       JSON object of API error
			2.           <AddUserDBError>         JSON object of API error
			3.           <ResponseCreatingError>  JSON object of API error

Login

"base/users/login" address:
	method                 POST
	required parameters:
		email                [email]            existing email of user
		password             [printable ASCII]  password which was used while creating
	return result:
		status 200:
			JSON object of user login confirm if request from "Android_app" and
			Set-Cookie with "session_id" key which is used for confidential actions
		status 400:
			1.           <RequestFormParseError>  JSON object of API error
			2.           <RequestFormDecodeError> JSON object of API error
			3.           <NoRequiredInfoError>    JSON object of API error
			4.           <RequestDataValidError>  JSON object of API error
			5.           <BadAuth>                JSON object of API error
		status 500:
			1.           <GetInfoDBError>         JSON object of API error
			2.           <ResponseCreatingError>  JSON object of API error
			3.           <SessionCreateError>     JSON object of API error

Logout

Cookie with tocken required to delete session. If cookie are not provided there is no effect.

"base/users/logout" address:
	method                 POST
	return result          status 200 always

Update information about user

Cookie with tocken required for this action.
If avatar_address is empty then avatar image will be deleted if exists.

"base/users/profile" address:
	method                 POST
	required parameters:
		first_name           [UTF letters]      first name of user
		last_name            [UTF letters]      last name of user
	allowed parameters:
		tel_number           [digits 1-9]       telephone number of user
		about                                   some additional information about user
		avatar_address       [existing address] address to existing user avatar
		images               [.JPEG or .png]    avatar image of user (if provided then all parameters must be in "multipart/form-data")
	return result:
		status 200           update succeed
		status 401:
			1.           <NoCookieError>          JSON object of API error
			2.           <BadCookieError>         JSON object of API error
		status 400:
			1.           <RequestFormParseError>  JSON object of API error
			2.           <RequestFormDecodeError> JSON object of API error
			3.           <NoRequiredInfoError>    JSON object of API error
			4.           <RequestDataValidError>  JSON object of API error
			5.           <ImageNoExistError>      JSON object of API error
		status 500:
			1.           <ImageCreateError>       JSON object of API error
			2.           <UpdateUserDBError>      JSON object of API error
			3.           <GetInfoDBError>         JSON object of API error

Get information about current logged user

Cookie required for this action.

"base/users/profile" address:
	method                 GET
	return result:
		status 200           JSON object of logged user
		status 401:
			1.           <NoCookieError>          JSON object of API error
			2.           <BadCookieError>         JSON object of API error
		status 500:
			1.           <GetInfoDBError>         JSON object of API error
			2.           <ResponseCreatingError>  JSON object of API error

Delete existing user

Cookie required for this action.

"base/users/profile" address:
	method                 DELETE
	return result:
		status 200           delete succeed
		status 401:
			1.           <NoCookieError>          JSON object of API error
			2.           <BadCookieError>         JSON object of API error
		status 500:
			1.           <GetInfoDBError>         JSON object of API error
			2.           <RemoveUserError>        JSON object of API error

Create new ad

Cookie required for this action.

"base/ads/new" address:
	method                 POST
	required parameters:
		title                                   title of ad
		city                                    city where ad is provided
		description_ad                          additional information about ad
	allowed parameters:
		price                [positive number]  price of ad
		country                                 country where ad is provided
		subway_station                          station where ad is provided
		images               [.JPEG or .png]    images of ad (if provided then all parameters must be in "multipart/form-data")
	return result:
		status 201           ad create confirm JSON object
		status 400:
			1.           <RequestFormParseError>  JSON object of API error
			2.           <RequestFormDecodeError> JSON object of API error
			3.           <NoRequiredInfoError>    JSON object of API error
			4.           <RequestDataValidError>  JSON object of API error
		status 401:
			1.           <NoCookieError>          JSON object of API error
			2.           <BadCookieError>         JSON object of API error
		status 500:
			1.           <ImageCreateError>       JSON object of API error
			2.           <CreateAdError>          JSON object of API error
			3.           <ResponseCreatingError>  JSON object of API error

Update existing ad

Cookie required for this action.
If parameter "ad_images" is empty then images will be deleted if exist.
If parameter "ad_images" is provided with existing addresses but content-type is
"multipart/data-form" and parameter "images" is not null then images will be appended
to existing.

"base/ads/edit/{id}" address:
	method               POST
	id                   must be a digit number
	required parameters:
		title                                   title of ad
		city                                    city where ad is provided
		description_ad                          additional information about ad
	allowed parameters:
		price                [positive number]     price of ad
		country                                    country where ad is provided
		subway_station                             station where ad is provided
		ad_images            [existing addresses]  array of existing addresses of ad's images
		images               [.JPEG or .png]       images of ad (if provided then all parameters must be in "multipart/form-data")
	return result:
		status 200           updating succeed
		status 400:
			1.           <RequestFormParseError>  JSON object of API error
			2.           <RequestFormDecodeError> JSON object of API error
			3.           <NoRequiredInfoError>    JSON object of API error
			4.           <RequestDataValidError>  JSON object of API error
			5.           <NoAdWithSuchIDError>    JSON object of API error
			6.           <ImageNoExistError>      JSON object of API error
		status 401:
			1.           <NoCookieError>          JSON object of API error
			2.           <BadCookieError>         JSON object of API error
		status 500:
			1.           <GetInfoDBError>         JSON object of API error
			2.           <ImageCreateError>       JSON object of API error
			3.           <UpdateAdError>          JSON object of API error

Delete existing ad

Cookie required for this action.

"base/ads/delete/{id}" address:
	method                 DELETE
	id                     must be a digit number
	return result:
		status 200           deleting succeed
		status 400:
			1.           <NoAdWithSuchIDError>    JSON object of API error
		status 401:
			1.           <NoCookieError>          JSON object of API error
			2.           <BadCookieError>         JSON object of API error
		status 500:
			1.           <GetInfoDBError>         JSON object of API error
			2.           <RemoveAdError>          JSON object of API error

Get images

"base/images/{filename}" address:
	method                 GET
	filename               must be existing image
	return result:
		status 200           image with such filename
		status 400           <NoSuchImageError> JSON object of API error
*/
package api

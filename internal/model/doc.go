// Copyright Dmitry Kargashin <dkargashin3@gmail.com>
// License can be found in LICENSE file.

/*
	Package model describes an internal interface of application.

	Data types

		User - a human that can sign up, sign in, logout, create ads and etc.
		Characteristics:
			id	                  unique identificator of user in database
			first name            first name of user
			last name	            last name of user
			email                 email of user that must be unique
			password              password of user that stored in database in hashed state
			telephone number      telephone number of user
			about                 some additional information about user

		Ad - main content of application
		Characteristics:
			id                    unique identificator of ad in database
			title                 title of ad
			price                 price of ad
			country               country where ad is situated
			city                  city where ad is situated
			subway station        station where ad is situated
			images                images of ad
			owner                 user which created this ad
			description           additional information about this ad
			creation time         time when this ad was created
*/
package model

# Search&Build project

[![pipeline status](https://bmstu.codes/developers34/SBWeb/badges/master/pipeline.svg)](https://bmstu.codes/developers34/SBWeb/commits/master)
[![coverage report](https://bmstu.codes/developers34/SBWeb/badges/master/coverage.svg)](https://bmstu.codes/developers34/SBWeb/commits/master)
[![Go Report Card](https://goreportcard.com/badge/bmstu.codes/developers34/SBWeb)](https://goreportcard.com/report/bmstu.codes/developers34/SBWeb)
[![License: LGPL v3](https://img.shields.io/badge/License-LGPL%20v3-blue.svg)](https://www.gnu.org/licenses/lgpl-3.0)
[![GoDoc](https://godoc.org/bmstu.codes/developers34/SBWeb?status.svg)](https://godoc.org/bmstu.codes/developers34/SBWeb)

## General overview of architecture

![Overview](/docs/GeneralOverview.png "Overview")

## Data structures

### User of service

He can publish ads. Consists of such information as:

* Unique identifier defined by API
* First name
* Last name
* Email
* Telephone number (can be omited)
* About (can be omited)
* Date and time of sign up

Example of JSON object that will be transmitted from API
to clients:

```json
{
    "id": 123456,
    "first_name": "Random",
    "last_name": "Valerka",
    "email": "valerka@example.com",
    "tel_num": "1-234-56-78",
    "about": "Some information about this man",
    "time_reg": "2012.10.1 15:34:41"
}
```

### Ad that can be published by users

* Unique identifier defined by API
* Title
* Price (can be omited)
* Country (can be omited)
* City
* Subway station (can be omited)
* Images (can be omited)
* Information about agent (just some fields from user structure)
* Description of service
* Creation date and time

Example of JSON object that will be transmitted from API
to clients:

```json
{
    "id": 1234,
    "title": "My awesome title",
    "price": 100500,
    "country": "Russia",
    "city": "Moscow",
    "subway_station": "Technopark",
    "images_url": ["ex.com/ad_id/1.png", "ex.com/ad_id/2.png"],
    "agent_info": {
        "id": 123456,
        "first_name": "Random",
        "last name": "Valerka",
        "email": "valerka@example.com",
        "tel_num": "1-234-56-78",
        "about": "Some information about this man",
        "time_reg": "2012.10.1 15:34:41"
    },
    "description": "it is awesome service with the best quality!",
    "time_cre": "2012.10.1 15:40:52"
}
```

### Error type

It will be sent to client in JSON format if something went wrong.

* Message
* Description
* Error Code

Example:

```json
{
    "message": "Can't create user",
    "description": "User with such email is already exists",
    "error": "UserEmailExists"
}
```

## API interface

### 1 stage

1 stage's task is to make simple CRUD interface so people can create, read, update, delete ads
without any authentification.
`root` is base domain of the API server.

`root/ads/new` is supposed to create new ads. Request on this URL must have method `POST` and must contain
several parameters in body that are simillar to JSON object of _ad_. URL will return http status corresponding to
result of creating *(need to define)*.

`root/ads/{id}` with method `GET` will return JSON object of _ad_. If there is no _ad_ with such id,
it will return http error *(need to define error code)*.

`root/ads` with method `GET` will return array of JSON objects of _ad_. URL receive two parameters:
`offset` defines the id of first _ad_; `count` defines the number of _ads_ that will be transmitted.
Default values are 0 and 10.

`root/ads/{id}` with method `POST` will update existing _ad_. URL will return http status corresponding to
result of updating *(need to define)*.

`root/ads/{id}` with method `DELETE` will delete existing _ad_. URL will return http status corresponding to
result of updating *(need to define)*.

`root/users/new` - create new user. Method must be the `POST`. Http body must contain parameters simillar to
JSON object of _user_. Email field must be unique among all users.

`root/users/{id}` - method `GET`. Show information about _user_ with such id. This URL will return JSON object of _user_.
With parameter `show_ads=true` added, URL will return array of JSON objects _ads_ of user with such id.

`root/users/{id}` - method `POST`. Update existing _user_.

`root/users/{id}` - method `DELETE`. Delete existing _user_.

### 2 stage. Autentification and authorization

#### Sign up

![SignUp](/docs/SignUp.PNG "SignUp")

1. Web-server, Android-application: receive data entered by user and validate it
2. Create request with method `POST` to `root/users/new`. Possible parameters:
    * __email__ (string) *email of user, must be unique*
    * __password__ (string) *password of user, can contain numbers and english letters*
    * __first_name__ (string) *first name of user, only english or russian letters allowed*
    * __last_name__ (string) *last name of user, only english or russian letters allowed*
    * __tel_num__ (string, *optional*) *telephone number of user in any format*
    * __about__ (string, *optional*) *information about user: what does he likes, useful skills etc.*
3. If there is no user with such email in API database, then record with user's data will be inserted in database. On this case API returns HTTP status `201 Created` and JSON object

    ```JSON with Comments
    {
        "id": 123,               // id of new user
        "ref": "root/users/123"  // URL to new user
    }
    ```

4. If user is already exists API will return HTTP status `409 Conflict` and JSON object of Error:

    ```json
    {
        "description": "User with such email is already exists",
        "message": "Can't create user",
        "error": "UserEmailExists"
    }
    ```

#### Sign in

![SignIn](/docs/SignIn.PNG "SignIn")

#### Access to actions requiring authorization

![AuthReq](/docs/AuthReq.PNG "AuthReq")

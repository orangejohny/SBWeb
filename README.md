# Search&Build project

[![pipeline status](https://bmstu.codes/developers34/SBWeb/badges/master/pipeline.svg)](https://bmstu.codes/developers34/SBWeb/commits/master)
[![coverage report](https://bmstu.codes/developers34/SBWeb/badges/master/coverage.svg)](https://bmstu.codes/developers34/SBWeb/commits/master)
[![Go Report Card](https://goreportcard.com/badge/bmstu.codes/developers34/SBWeb)](https://goreportcard.com/report/bmstu.codes/developers34/SBWeb)
[![License: LGPL v3](https://img.shields.io/badge/License-LGPL%20v3-blue.svg)](https://www.gnu.org/licenses/lgpl-3.0)
[![GoDoc](https://godoc.org/github.com/orangejohny/SBWeb?status.svg)](https://godoc.org/github.com/orangejohny/SBWeb)

## Links

* [wiki](https://bmstu.codes/developers34/SBWeb/wikis/home)
* [GoDoc](https://godoc.org/github.com/orangejohny/SBWeb)
* [Coverage](https://developers34.pages.bmstu.codes/SBWeb)
* [Mirror on GitHub](https://github.com/orangejohny/SBWeb)
* [Heroku App](https://search-build.herokuapp.com)

## Wiki

**Check the [Wiki](https://bmstu.codes/developers34/SBWeb/wikis/home) page for more information!**

## Overview

**Search&Build** is a service that allows you to solve a wide spectre of building tasks easilly.  
It's online marketplace of various building services: you can hire different specialists,  
buy some services that connected to construction industry, you can even offer your own services  
and use **Search&Build** to make your business more effective due to the wide reach of the target audience.

## Install

```bash
go get -u bmstu.codes/developers34/SBWeb
```

More information is on [wiki](https://bmstu.codes/developers34/SBWeb/wikis/Install-and-usage).

## Interface

Information about interface is [here](https://godoc.org/github.com/orangejohny/SBWeb/pkg/api).

Allowed addresses:
* /ads                    `GET`
* /ads/{id}               `GET`
* /users/{id}             `GET`
* /users/new              `POST`
* /users/login            `POST`
* /users/logout           `POST`
* /users/profile          `GET`
* /users/profile          `POST`
* /users/profile          `DELETE`
* /ads/new                `POST`
* /ads/edit/{id}          `POST`
* /ads/delete/{id}        `DELETE`
* /images/{filename}      `GET`
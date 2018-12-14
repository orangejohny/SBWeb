// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

package db_test

import (
	"database/sql"
	"testing"

	"gopkg.in/guregu/null.v3/zero"

	"bmstu.codes/developers34/SBWeb/pkg/model"

	"bmstu.codes/developers34/SBWeb/pkg/db"
)

func TestInit(t *testing.T) {
	cfg := db.Config{
		DBAddress:    "postgresql://runner:@postgres/data?sslmode=disable",
		MaxOpenConns: 10,
	}

	database, _ := sql.Open("postgres", "postgresql://runner:@postgres/data?sslmode=disable")
	database.Exec(`
	CREATE TABLE IF NOT EXISTS users
(
    id                SERIAL      PRIMARY KEY,
    first_name        varchar(80) NOT NULL,
    last_name         varchar(80) NOT NULL,
    email             varchar(80) UNIQUE NOT NULL,
    password_hash     text        NOT NULL,
    telephone         varchar(80),
    about             text,
    avatar_address    text,
    reg_time          timestamp   DEFAULT CURRENT_TIMESTAMP NOT NULL
);`)
	database.Exec(`
	CREATE TABLE IF NOT EXISTS ads
(
    id             SERIAL       PRIMARY KEY,
    title          varchar(80)  NOT NULL,
    price          integer      CONSTRAINT positive_price CHECK (price > 0),
    country        varchar(80),
    city           varchar(80),
    subway_station varchar(80),
    ad_images      varchar(256)[],
    -- when deleting user we should delete his ads
    owner_ad       integer      REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    description_ad text,
    creation_time  timestamp    DEFAULT CURRENT_TIMESTAMP NOT NULL
);`)

	database.Close()

	_, err := db.InitConnDB(cfg)
	if err != nil {
		t.Error("Unexpected error", err.Error())
	}

	cfg = db.Config{
		DBAddress:    "postgresql://rustgres/data?sslmode=disable",
		MaxOpenConns: 10,
	}

	_, err = db.InitConnDB(cfg)
	if err == nil {
		t.Error("Expected error")
	}

	cfg = db.Config{
		DBAddress:    "postgresql://runner:@postgres/data?sslmode=disable",
		MaxOpenConns: 10,
	}

	database, _ = sql.Open("postgres", "postgresql://runner:@postgres/data?sslmode=disable")
	database.Exec("DROP TABLE ads")
	database.Close()

	_, err = db.InitConnDB(cfg)
	if err == nil {
		t.Error("Expected error")
	}
}

func TestInterface(t *testing.T) {
	cfg := db.Config{
		DBAddress:    "postgresql://runner:@postgres/data?sslmode=disable",
		MaxOpenConns: 10,
	}

	database, _ := sql.Open("postgres", "postgresql://runner:@postgres/data?sslmode=disable")
	database.Exec(`
	CREATE TABLE IF NOT EXISTS users
(
    id                SERIAL      PRIMARY KEY,
    first_name        varchar(80) NOT NULL,
    last_name         varchar(80) NOT NULL,
    email             varchar(80) UNIQUE NOT NULL,
    password_hash     text        NOT NULL,
    telephone         varchar(80),
    about             text,
    avatar_address    text,
    reg_time          timestamp   DEFAULT CURRENT_TIMESTAMP NOT NULL
);`)
	database.Exec(`
	CREATE TABLE IF NOT EXISTS ads
(
    id             SERIAL       PRIMARY KEY,
    title          varchar(80)  NOT NULL,
    price          integer      CONSTRAINT positive_price CHECK (price > 0),
    country        varchar(80),
    city           varchar(80),
    subway_station varchar(80),
    ad_images      varchar(256)[],
    -- when deleting user we should delete his ads
    owner_ad       integer      REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    description_ad text,
    creation_time  timestamp    DEFAULT CURRENT_TIMESTAMP NOT NULL
);`)

	database.Exec(`INSERT INTO users
	(first_name, last_name, email, password_hash)
	VALUES
	('Ivan', 'Ivanov', 'ivan@gmail.com', '123456')`)

	database.Exec(`INSERT INTO ads
	(title, owner_ad, description_ad, price, city)
	VALUES
	('Building', 1, 'Some description', '100500', 'Moscow')`)

	database.Exec(`INSERT INTO ads
	(title, owner_ad, description_ad, price, city)
	VALUES
	('Building that can be built', 1, 'Some description', 100500, 'Moscow')`)

	database.Close()

	h, err := db.InitConnDB(cfg)
	if err != nil {
		t.Error("Unexpected error", err.Error())
	}

	user := model.User{
		ID:        1,
		FirstName: "Ivan",
		LastName:  "Ivanov",
		Email:     "ivan@gmail.com",
	}

	userNew := model.User{
		FirstName: "Alex",
		LastName:  "Ivanov",
		Email:     "alex@gmail.com",
		Password:  "123455",
	}

	adNew := model.AdItem{
		Title:       "New ad",
		UserID:      1,
		Description: "Gooooood",
		Price:       zero.NewInt(500, true),
		City:        "New York",
	}

	ad1 := model.AdItem{
		ID:          1,
		Title:       "Building",
		UserID:      1,
		User:        user,
		Description: "Some description",
		Price:       zero.NewInt(100500, true),
		City:        "Moscow",
	}

	ad2 := model.AdItem{
		ID:          2,
		Title:       "Building that can be built",
		UserID:      1,
		User:        user,
		Description: "Some description",
		Price:       zero.NewInt(100500, true),
		City:        "Moscow",
	}

	ads, err := h.GetAds(&model.SearchParams{
		Limit:  15,
		Offset: 0,
	})
	if err != nil {
		t.Error("Unexpected error", err.Error())
	} else if len(ads) != 2 {
		t.Error("Unexpected len", len(ads))
	} else if ads[0].Title != ad1.Title || ads[1].Title != ad2.Title {
		t.Error("Expected equal ads")
	}

	ads, err = h.GetAds(&model.SearchParams{
		Limit:  15,
		Query:  "that",
		Offset: 0,
	})
	if err != nil {
		t.Error("Unexpected error", err.Error())
	} else if len(ads) != 1 {
		t.Error("Unexpected len", len(ads))
	} else if ads[0].Title != ad2.Title {
		t.Error("Expected equal ads")
	}

	ads, err = h.GetAdsOfUser(1)
	if err != nil {
		t.Error("Unexpected error", err.Error())
	} else if len(ads) != 2 {
		t.Error("Unexpected len", len(ads))
	} else if ads[0].Title != ad1.Title || ads[1].Title != ad2.Title {
		t.Error("Expected equal ads")
	}

	ad, err := h.GetAd(1)
	if err != nil {
		t.Error("Unexpected error", err.Error())
	} else if ad.Title != ad1.Title {
		t.Error("Expected equal ads")
	}

	ad, _ = h.GetAd(15)
	if ad.ID != -1 {
		t.Error("Expected ID = -1")
	}

	u, err := h.GetUserWithID(1)
	if err != nil {
		t.Error("Unexpected error", err.Error())
	} else if u.FirstName != user.FirstName {
		t.Error("Expected equal users")
	}

	u, _ = h.GetUserWithID(15)
	if u.ID != -1 {
		t.Error("Expected ID = -1")
	}

	u, err = h.GetUserWithEmail("ivan@gmail.com")
	if err != nil {
		t.Error("Unexpected error", err.Error())
	} else if u.FirstName != user.FirstName {
		t.Error("Expected equal users")
	}

	u, _ = h.GetUserWithEmail("feffr2C")
	if u.ID != -1 {
		t.Error("Expected ID = -1")
	}

	id, _ := h.NewUser(&user)
	if id != -1 {
		t.Error("Expected ID = -1 got = ", id)
	}

	_, err = h.NewUser(&userNew)
	if err != nil {
		t.Error("Unexpected error", err.Error())
	}

	id, err = h.NewAd(&adNew)
	if err != nil {
		t.Error("Unexpected error", err.Error())
	} else if id != 3 {
		t.Error("Expected id = 3 got = ", id)
	}

	id, err = h.EditUser(&user)
	if err != nil {
		t.Error("Unexpected error", err.Error())
	} else if id != 1 {
		t.Error("Expected id = 1 got = ", id)
	}

	id, err = h.EditAd(&ad1)
	if err != nil {
		t.Error("Unexpected error", err.Error())
	} else if id != 1 {
		t.Error("Expected id = 1 got = ", id)
	}

	id, err = h.RemoveAd(1)
	if err != nil {
		t.Error("Unexpected error", err.Error())
	} else if id != 1 {
		t.Error("Expected id = 1 got = ", id)
	}

	id, err = h.RemoveUser(1)
	if err != nil {
		t.Error("Unexpected error", err.Error())
	} else if id != 1 {
		t.Error("Expected id = 1 got = ", id)
	}
}

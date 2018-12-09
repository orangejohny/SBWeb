// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

package db_test

import (
	"database/sql"
	"reflect"
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
	('Bulding', 1, 'Some description', '100500', 'Moscow')`)

	database.Exec(`INSERT INTO ads
	(title, owner_ad, description_ad, price, city)
	VALUES
	('Bulding that can be built', 1, 'Some description', 100500, 'Moscow')`)

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
	} else if !reflect.DeepEqual(*ads[0], ad1) || !reflect.DeepEqual(*ads[1], ad2) {
		t.Error("Expected equal ads.Received:", *ads[0], "and\n", *ads[1])
	}
}

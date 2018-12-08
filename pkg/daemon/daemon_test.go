// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

package daemon_test

import (
	"database/sql"
	"testing"

	"bmstu.codes/developers34/SBWeb/pkg/api"
	"bmstu.codes/developers34/SBWeb/pkg/daemon"
	"bmstu.codes/developers34/SBWeb/pkg/db"
	"bmstu.codes/developers34/SBWeb/pkg/sessionmanager"
)

func TestRunService(t *testing.T) {
	cfg := &daemon.Config{
		DB: db.Config{
			DBAddress:    "postgresql://runner:@postgres/data?sslmode=disable",
			MaxOpenConns: 10,
		},
		SM: sessionmanager.Config{
			DBAddress:      "redis://redis:6379/0",
			TockenLength:   32,
			ExpirationTime: 86400,
		},
		API: api.Config{
			Address:      ":54000",
			ReadTimeout:  "10s",
			WriteTimeout: "10s",
			IdleTimeout:  "10s",
		},
	}

	db, _ := sql.Open("postgres", "postgresql://runner:@postgres/data?sslmode=disable")
	db.Exec(`
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
	db.Exec(`
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

	db.Close()

	err := daemon.RunService(cfg)
	if err != nil {
		t.Error("Unexpected error: ", err.Error())
	}
}

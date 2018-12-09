// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

package model

import (
	"time"

	"gopkg.in/guregu/null.v3/zero"
)

// AdItem struct describes ad which users are supposed to create, watch and etc.
type AdItem struct {
	ID            int64       `db:"idad" json:"id" schema:"id,optional" valid:"-"`
	Title         string      `db:"title" json:"title" schema:"title,optional" valid:",optional"` // required in DB
	Price         zero.Int    `db:"price" json:"price,omitempty" schema:"price,optional" valid:"-"`
	Country       zero.String `db:"country" json:"country,omitempty" schema:"country,optional" valid:"-"`                      // consists of printable ASCII
	City          string      `db:"city" json:"city,omitempty" schema:"city,optional" valid:",optional"`                       // required in DB
	SubwayStation zero.String `db:"subway_station" json:"subway_station,omitempty" schema:"subway_station,optional" valid:"-"` // consists of printable ASCII
	AdImages      []string    `db:"-" json:"ad_images,omitempty" schema:"ad_images,optional" valid:"-"`
	AdImagesStr   zero.String `db:"ad_images" json:"-" schema:"-" valid:"-"` // for database
	UserID        int64       `db:"owner_ad" json:"-" schema:"-" valid:"-"`  // for database
	User          `json:"owner_ad" schema:"-" valid:"-"`
	Description   string    `db:"description_ad" json:"description_ad" schema:"description_ad,optional" valid:",optional"` // requiered in DB
	CreationTime  time.Time `db:"creation_time" json:"creation_time" schema:"-" valid:"-"`
}

// TODO country, city, subway station should be UTF letters with some characters
// description should be valid UTF-8

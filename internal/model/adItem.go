package model

import (
	"time"

	"gopkg.in/guregu/null.v3/zero"
)

// AdItem struct describes ad that users supposed to add
type AdItem struct {
	ID            int64       `db:"id_" json:"id" schema:"id,optional" valid:"-"`
	Title         string      `db:"title" json:"title" schema:"title,optional" valid:"printableascii,optional"` // required in DB
	Price         zero.Int    `db:"price" json:"price,omitempty" schema:"price,optional" valid:"-"`
	Country       zero.String `db:"country" json:"country,omitempty" schema:"country,optional" valid:"-"`                      // consists of utf letter
	City          string      `db:"city" json:"city,omitempty" schema:"city,optional" valid:"utfletter,optional"`              // required in DB
	SubwayStation zero.String `db:"subway_station" json:"subway_station,omitempty" schema:"subway_station,optional" valid:"-"` // consists of utf letter
	ImagesFolder  zero.String `db:"images_folder" json:"images_folder,omitempty" schema:"-" valid:"-"`                         // will be implemented (maybe)
	UserID        int64       `db:"owner_ad" json:"-" schema:"-" valid:"-"`
	User          `json:"owner_ad" schema:"-" valid:"-"`
	Description   string    `db:"description_ad" json:"description_ad" schema:"description_ad,optional" valid:"ascii,optional"` // requiered in DB
	CreationTime  time.Time `db:"creation_time" json:"creation_time" schema:"-" valid:"-"`
}

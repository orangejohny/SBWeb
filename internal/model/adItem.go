package model

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

// AdItem struct describes ad that users supposed to add
type AdItem struct {
	ID            int64       `db:"id_" json:"id" schema:"id,optional" valid:"-"`
	Title         string      `db:"title" json:"title" schema:"title,required" valid:",required"`
	Price         null.Int    `db:"price" json:"price,omitempty" schema:"price,optional" valid:"-"`
<<<<<<< HEAD
	Country       null.String `db:"country" json:"country,omitempty" schema:"country,optional" valid:"alpha,optional"`
	City          string      `db:"city" json:"city,omitempty" schema:"city,required" valid:"alpha,optional"`
	SubwayStation null.String `db:"subway_station" json:"subway_station,omitempty" schema:"subway_station,optional" valid:"alpha,optional"`
=======
	Country       null.String `db:"country" json:"country,omitempty" schema:"country,optional" valid:"-"`
	City          null.String `db:"city" json:"city,omitempty" schema:"city,required" valid:"-"`
	SubwayStation null.String `db:"subway_station" json:"subway_station,omitempty" schema:"subway_station,optional" valid:"-"`
>>>>>>> 025df3e48dd0859b2eda53f467bfabe9ed52dad9
	ImagesFolder  null.String `db:"images_folder" json:"images_folder,omitempty" schema:"-" valid:"-"`
	UserID        int64       `db:"owner_ad" json:"-" schema:"-" valid:"-"`
	User          `json:"owner_ad" schema:"-" valid:"-"`
	Description   string    `db:"description_ad" json:"description_ad" schema:"description_ad,optional" valid:",required"`
	CreationTime  time.Time `db:"creation_time" json:"creation_time" schema:"-" valid:"-"`
}

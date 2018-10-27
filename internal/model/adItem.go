package model

import (
	"time"
)

// TODO add tags for JSON, validator, schema

// AdItem struct describes ad that users supposed to add
type AdItem struct {
	ID            int       `db:"id" json:"id"`
	Title         string    `db:"title" json:"title"`
	Price         int       `db:"price" json:"price,omitempty"`
	Country       string    `db:"country" json:"country,omitempty"`
	City          string    `db:"city" json:"city,omitempty"`
	SubwayStation string    `db:"subway_station" json:"subway_station,omitempty"`
	ImagesFolder  string    `db:"images_folder" json:"images_folder,omitempty"`
	Owner         User      `db:"owner_ad" json:"owner_ad"`
	Description   string    `db:"description_ad" json:"description_ad"`
	CreationTime  time.Time `db:"creation_time" json:"creation_time"`
}

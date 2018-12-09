// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"log"
	"strings"

	"bmstu.codes/developers34/SBWeb/pkg/model"
)

const (
	notUniqueEmail = `pq: duplicate key value violates unique constraint "users_email_key"`
)

// prepareStateents preapares SQL statements for interaction with postgres database.
func (h *Handler) prepareStatements() (err error) {
	if h.ReadAds, err = h.DB.PrepareNamed( // return list of ads
		`SELECT
		 ads.id "idad", title, description_ad, price, country, city, subway_station, array_to_string(ad_images,',') "ad_images", creation_time, owner_ad,
		 users.id, first_name, last_name, email, telephone, about, reg_time, avatar_address
		 FROM
		 ads
		 INNER JOIN
		 users 
		 ON
		 users.id = ads.owner_ad
		 LIMIT :limit OFFSET :offset`,
	); err != nil {
		log.Println(err.Error())

		return err
	}

	if h.SearchAds, err = h.DB.PrepareNamed(
		`SELECT
		ads.id "idad", title, description_ad, price, country, city, subway_station, array_to_string(ad_images,',') "ad_images", creation_time, owner_ad,
		users.id, first_name, last_name, email, telephone, about, reg_time, avatar_address
		FROM
		ads
		INNER JOIN
		users 
		ON
		users.id = ads.owner_ad
		WHERE ads.title ILIKE '%' || :query || '%'
		LIMIT :limit OFFSET :offset`,
	); err != nil {
		log.Println(err.Error())

		return err
	}

	if h.ReadAdsOfUser, err = h.DB.Preparex( // return list of ads of such user
		`SELECT
		 ads.id "idad", title, description_ad, price, country, city, subway_station, array_to_string(ad_images,',') "ad_images", creation_time, owner_ad,
		 users.id, first_name, last_name, email, telephone, about, reg_time, avatar_address
		 FROM
		 ads
		 INNER JOIN
		 users 
		 ON
		 users.id = ads.owner_ad
		 AND ads.owner_ad = $1`,
	); err != nil {
		log.Println(err.Error())

		return err
	}

	if h.ReadAd, err = h.DB.Preparex( // return ad with such id
		`SELECT
		ads.id "idad", title, description_ad, price, country, city, subway_station, array_to_string(ad_images,',') "ad_images", creation_time, owner_ad,
		users.id, first_name, last_name, email, telephone, about, reg_time, avatar_address
		FROM
		ads
		INNER JOIN
		users
		ON
		users.id = ads.owner_ad AND ads.id = $1`,
	); err != nil {
		log.Println(err.Error())

		return err
	}

	if h.ReadUserWithID, err = h.DB.Preparex( // return user with such id
		"SELECT id, first_name, last_name, email, telephone, about, reg_time, avatar_address FROM users WHERE id=$1",
	); err != nil {
		log.Println(err.Error())

		return err
	}

	if h.ReadUserWithEmail, err = h.DB.Preparex( // return user with such email
		"SELECT id, first_name, last_name, email, telephone, about, reg_time, password_hash, avatar_address FROM users WHERE email=$1",
	); err != nil {
		log.Println(err.Error())

		return err
	}

	if h.CreateUser, err = h.DB.PrepareNamed( // create new user
		`INSERT INTO users
			(first_name, last_name, email, password_hash, telephone, about, avatar_address)
			VALUES
			(:first_name, :last_name, :email, :password_hash, :telephone, :about, :avatar_address)
			RETURNING id`,
	); err != nil {
		log.Println(err.Error())

		return err
	}

	if h.CreateAd, err = h.DB.PrepareNamed( // create new ad
		`INSERT INTO ads
			(title, owner_ad, description_ad, price, country, city, subway_station, ad_images)
			VALUES
			(:title, :owner_ad, :description_ad, :price, :country, :city, :subway_station, string_to_array(:ad_images, ','))
			RETURNING id`,
	); err != nil {
		log.Println(err.Error())

		return err
	}

	if h.UpdateUser, err = h.DB.PrepareNamed( // update user
		`UPDATE users SET
			first_name=:first_name,
			last_name=:last_name,
			telephone=:telephone,
			about=:about,
			avatar_address=:avatar_address
			WHERE id=:id`,
	); err != nil {
		log.Println(err.Error())

		return err
	}

	if h.UpdateAd, err = h.DB.PrepareNamed( // update ad
		`UPDATE ads SET
			title=:title,
			description_ad=:description_ad,
			price=:price,
			country=:country,
			city=:city,
			subway_station=:subway_station,
			ad_images=string_to_array(:ad_images, ',')
			WHERE id=:idad`,
	); err != nil {
		log.Println(err.Error())

		return err
	}

	if h.DeleteUser, err = h.DB.Preparex( // delete user
		`DELETE FROM users WHERE id=$1`,
	); err != nil {
		log.Println(err.Error())

		return err
	}

	if h.DeleteAd, err = h.DB.Preparex( // delete ad
		`DELETE FROM ads WHERE id=$1`,
	); err != nil {
		log.Println(err.Error())

		return err
	}

	return nil
}

// GetAds returns slice of model.AdItem from database based on incoming filters.
func (h *Handler) GetAds(sp *model.SearchParams) ([]*model.AdItem, error) {
	ads := make([]*model.AdItem, 0)
	var err error
	if sp.Query == "" {
		err = h.ReadAds.Select(&ads, sp)
	} else {
		err = h.SearchAds.Select(&ads, sp)
	}

	return ads, err
}

// GetAdsOfUser returns slice of model.AdItem with such user from database.
func (h *Handler) GetAdsOfUser(userID int64) ([]*model.AdItem, error) {
	ads := make([]*model.AdItem, 0)
	err := h.ReadAdsOfUser.Select(&ads, userID)
	return ads, err
}

// GetAd returns model.AdItem struct with such ID.
func (h *Handler) GetAd(adID int64) (*model.AdItem, error) {
	ad := &model.AdItem{}
	err := h.ReadAd.Get(ad, adID)
	ad.AdImages = strings.Split(ad.AdImagesStr.String, ",")
	if err == sql.ErrNoRows {
		ad.ID = -1
	}
	return ad, err
}

// GetUserWithID returns model.User struct with such ID.
func (h *Handler) GetUserWithID(userID int64) (*model.User, error) {
	user := &model.User{}
	err := h.ReadUserWithID.Get(user, userID)
	if err == sql.ErrNoRows { // is 'false' possible?
		user.ID = -1
	}
	return user, err
}

// GetUserWithEmail returns model.User struct with such email.
func (h *Handler) GetUserWithEmail(email string) (*model.User, error) {
	user := &model.User{}
	err := h.ReadUserWithEmail.Get(user, email)
	if err == sql.ErrNoRows {
		user.ID = -1
	}
	return user, err
}

// NewUser adds new User to database if it is possible.
func (h *Handler) NewUser(user *model.User) (int64, error) {
	var lastInserted int64

	err := h.CreateUser.Get(&lastInserted, user)
	if err != nil && err.Error() == notUniqueEmail {
		lastInserted = -1
	}

	return lastInserted, err
}

// NewAd creates a new row in "ads" table in database.
func (h *Handler) NewAd(ad *model.AdItem) (int64, error) {
	var lastInserted int64
	ad.AdImagesStr.SetValid(strings.Join(ad.AdImages, ","))
	err := h.CreateAd.Get(&lastInserted, ad)

	return lastInserted, err
}

// EditUser updates User with ID provided from function argument.
func (h *Handler) EditUser(user *model.User) (int64, error) {
	res, err := h.UpdateUser.Exec(user)
	if err != nil {
		return -1, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}

	return affected, nil
}

// EditAd updates information about ad with ID provided from function argument.
func (h *Handler) EditAd(ad *model.AdItem) (int64, error) {
	ad.AdImagesStr.SetValid(strings.Join(ad.AdImages, ","))

	res, err := h.UpdateAd.Exec(ad)
	if err != nil {
		return -1, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}

	return affected, nil
}

// RemoveUser deletes user with such ID from database.
func (h *Handler) RemoveUser(userID int64) (int64, error) {
	res, err := h.DeleteUser.Exec(userID)
	if err != nil {
		return -1, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}

	return affected, nil
}

// RemoveAd deletes ad with such ID from database.
func (h *Handler) RemoveAd(adID int64) (int64, error) {
	res, err := h.DeleteAd.Exec(adID)
	if err != nil {
		return -1, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}

	return affected, nil
}

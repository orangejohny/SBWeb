package db

import (
	"log"

	"github.com/orangejohny/SBWeb/internal/model"
)

// prepareStateents just preapares SQL requests
func (h *Handler) prepareStatements() (err error) {
	if h.ReadAds, err = h.DB.Preparex( // return list of ads
		"SELECT * FROM ads LIMIT $1 OFFSET $2",
	); err != nil {
		log.Println(err.Error())

		return err
	}

	if h.ReadAd, err = h.DB.Preparex( // return ad with such id
		"SELECT * FROM ads WHERE id=$1",
	); err != nil {
		log.Println(err.Error())

		return err
	}

	if h.ReadUserWithID, err = h.DB.Preparex( // return user with such id
		"SELECT * FROM users WHERE id=$1",
	); err != nil {
		log.Println(err.Error())

		return err
	}

	if h.ReadUserWithEmail, err = h.DB.Preparex( // return user with such email
		"SELECT * FROM users WHERE email=$1",
	); err != nil {
		log.Println(err.Error())

		return err
	}

	if h.CreateUser, err = h.DB.PrepareNamed( // create new user
		`INSERT INTO users
			(first_name, last_name, email, password_hash, telephone, about)
			VALUES
			(:first_name, :last_name, :email, :password_hash, :telephone, :about)
			RETURNING id`,
	); err != nil {
		log.Println(err.Error())

		return err
	}

	if h.CreateAd, err = h.DB.PrepareNamed( // create new ad
		`INSERT INTO ads
			(title, owner_ad, description_ad, price, country, city, subway_station, images_folder)
			VALUES
			(:title, :owner_ad, :description_ad, :price, :country, :city, :subway_station, :images_folder)`,
	); err != nil {
		log.Println(err.Error())

		return err
	}

	if h.UpdateUser, err = h.DB.PrepareNamed( // update user
		`UPDATE users SET
			first_name=:first_name,
			last_name=:last_name,
			telephone=:telephone,
			about=:about
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
			subway_station=:subway_station
			WHERE id=:id`,
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

// GetAds returns slice of AdItem from database
func (h *Handler) GetAds(limit int, offset int) ([]*model.AdItem, error) {
	ads := make([]*model.AdItem, 0)
	err := h.ReadAds.Select(&ads, limit, offset) // will sqlx manage with foreign keys?
	return ads, err
}

// GetAd returns AdItem struct with such ID
func (h *Handler) GetAd(adID int64) (*model.AdItem, error) {
	ad := &model.AdItem{}
	err := h.ReadAd.Get(ad, adID) // will sqlx manage with foreign keys?
	return ad, err
}

// GetUserWithID returns User struct with such ID
func (h *Handler) GetUserWithID(userID int64) (*model.User, error) {
	user := &model.User{}
	err := h.ReadUserWithID.Get(user, userID)
	return user, err
}

// GetUserWithEmail returns User struct with such email
func (h *Handler) GetUserWithEmail(email string) (*model.User, error) {
	user := &model.User{}
	err := h.ReadUserWithEmail.Get(user, email)
	return user, err
}

// NewUser adds new User to database if it is possible
func (h *Handler) NewUser(user *model.User) (int64, error) {
	var lastInserted int64

	err := h.CreateUser.Get(&lastInserted, user)
	if err != nil {
		return -1, err
	}

	return lastInserted, nil
}

// NewAd adds AdItem to database
func (h *Handler) NewAd(ad *model.AdItem) (int64, error) {
	var lastInserted int64

	err := h.CreateAd.Get(&lastInserted, ad)
	if err != nil {
		return -1, err
	}

	return lastInserted, nil
}

// EditUser updates User with ID provided from function argument
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

// EditAd updates information about ad with ID provided from function argument
func (h *Handler) EditAd(ad *model.AdItem) (int64, error) {
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

// RemoveUser deletes user with such ID from database
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

// RemoveAd deletes ad with such ID from database
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

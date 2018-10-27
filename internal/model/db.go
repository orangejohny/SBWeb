package model

// db describes interface of database needed by API
// to communicate with it
type db interface {
	GetAd(adID int) (*AdItem, error)
	GetAds(limit int, offset int) ([]*AdItem, error)
	GetUserWithID(userID int) (*User, error)
	GetUserWithEmail(email string) (*User, error) // no available for now
	NewUser(user *User) (int64, error)
	NewAd(ad *AdItem) (int64, error)
	EditUser(user *User) (int64, error)
	EditAd(ad *AdItem) (int64, error)
	RemoveUser(userID int) (int64, error)
	RemoveAd(adID int) (int64, error)
}

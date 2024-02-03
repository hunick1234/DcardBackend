package storage

import "github.com/hunick1234/DcardBackend/model"

type Database interface {
	StoreAd(ad *model.AD) error
	GetAd() (*model.AD, error)
}

type Cache interface {
	SetAd(ad *model.AD) error
	GetAd() (*model.AD, error)
}

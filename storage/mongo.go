package storage

import (
	"github.com/hunick1234/DcardBackend/model"
)

type MongoDB struct {
}

func NewMongoDB() *MongoDB {
	return &MongoDB{}
}

func (m *MongoDB) StoreAd(ad *model.AD) error {
	return nil
}

func (m *MongoDB) GetAd() (*model.AD, error) {
	return nil, nil
}

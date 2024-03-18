package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hunick1234/DcardBackend/dto"
	"github.com/hunick1234/DcardBackend/model/ad"
	"github.com/hunick1234/DcardBackend/myhttp"
	"github.com/hunick1234/DcardBackend/service"
)

type AdParama interface {
	ad.AD | ad.AdQuery
}

func GetAd(svc service.AdService, dto dto.Request, res *myhttp.Response) error {
	query := dto.GetRequestAdQuery()

	ad, err := svc.FindByFilter(context.Background(), &query)
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(ad)
	if err != nil {
		return err
	}
	res.Message = bytes
	res.StausCode = http.StatusOK
	return nil
}

func PostAd(service service.AdService, dto dto.Request, res *myhttp.Response) error {
	ad := dto.GetRequestAd()

	err := service.Store(context.Background(), &ad)
	if err != nil {
		return err
	}

	res.Message = []byte("success")
	res.StausCode = http.StatusOK
	return nil
}

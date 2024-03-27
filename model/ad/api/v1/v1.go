package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hunick1234/DcardBackend/model/ad"
	"github.com/hunick1234/DcardBackend/service"
	"github.com/hunick1234/DcardBackend/types"
)

type AdParama interface {
	ad.AD | ad.AdQuery
}

func GetAd(svc service.AdService, adCtx *types.AdControllerCtx) error {
	
	query := adCtx.R.GetRequestAdQuery()

	ad, err := svc.FindByFilter(context.Background(), &query)
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(ad)
	if err != nil {
		return err
	}
	adCtx.W.Message = bytes
	adCtx.W.StausCode = http.StatusOK
	return nil
}

func PostAd(service service.AdService, adCtx *types.AdControllerCtx) error {
	ad := adCtx.R.GetRequestAd()

	err := service.Store(context.Background(), &ad)
	if err != nil {
		return err
	}

	adCtx.W.Message = []byte("success")
	adCtx.W.StausCode = http.StatusOK
	return nil
}

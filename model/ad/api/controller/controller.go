package controller

import (
	"github.com/hunick1234/DcardBackend/service"
	"github.com/hunick1234/DcardBackend/types"
)

type APIController interface {
	InitEvent(*types.AdControllerCtx, service.AdService) error
	BeforeAPIEvent(*types.AdControllerCtx, service.AdService) error
	AfterAPIEvent(*types.AdControllerCtx, service.AdService) error
}

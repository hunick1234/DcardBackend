package types

import (
	"context"

	"github.com/hunick1234/DcardBackend/dto"
	"github.com/hunick1234/DcardBackend/myhttp"
)

type AdControllerCtx struct {
	Ctx context.Context
	R   *dto.Request
	W   *myhttp.Response
}

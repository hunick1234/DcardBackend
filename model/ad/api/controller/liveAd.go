package controller

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/hunick1234/DcardBackend/service"
	"github.com/hunick1234/DcardBackend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LiveAd struct {
	LiveAdCount int
	StartAt     int64
	lock        *sync.Mutex
}

func NewLiveAd() *LiveAd {
	return &LiveAd{
		LiveAdCount: 0,
		StartAt:     0,
		lock:        &sync.Mutex{},
	}
}

func (l *LiveAd) Pipeline() mongo.Pipeline {
	return mongo.Pipeline{
		{{"$match", bson.D{
			{"start_timestamp", bson.D{{"$lte", l.StartAt}}},
			{"end_timestamp", bson.D{{"$gte", l.StartAt}}},
		}}},
		{{"$group", bson.D{
			{"_id", nil},
			{"live_ad", bson.D{{"$count", bson.D{}}}}}},
		}}
}

func (l *LiveAd) BeforeAPIEvent(adCtx *types.AdControllerCtx, srv service.AdService) error {

	l.lock.Lock()
	defer l.lock.Unlock()
	ctx := adCtx.Ctx
	ad := adCtx.R.GetRequestAd()
	//remember set start_at time
	l.StartAt = ad.StartTimestamp
	err := l.event(ctx, srv)
	if err != nil {
		log.Println("live", err)
		adCtx.W.StausCode = http.StatusInternalServerError
		adCtx.Err = err
		return nil
	}
	if l.LiveAdCount >= 1000 {
		log.Println("live ad limit of 1000 data entries")
		adCtx.W.StausCode = http.StatusBadRequest
		adCtx.Err = errors.New("live ad limit of 1000 data entries")
		return nil
	}

	return nil
}

func (l *LiveAd) AfterAPIEvent(adCtx *types.AdControllerCtx, srv service.AdService) error {
	return nil
}

func (l *LiveAd) InitEvent(adCtx *types.AdControllerCtx, srv service.AdService) error {
	l.StartAt = 0
	l.LiveAdCount = 0
	return nil
}

func (l *LiveAd) SetStartAt(startAt int64) {
	l.StartAt = startAt
}

func (l *LiveAd) event(ctx context.Context, srv service.AdService) error {
	var aggResult struct {
		LiveAd int `bson:"live_ad"`
	}
	// Translate the result of the aggregation to a Go struct
	err := srv.Aggregate(context.TODO(), l, &aggResult)
	if err != nil {
		return err
	}
	l.LiveAdCount = aggResult.LiveAd
	log.Println("LiveAdCount", l.LiveAdCount)
	return nil
}

package controller

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/hunick1234/DcardBackend/service"
	"github.com/hunick1234/DcardBackend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DailyAd struct {
	DailyAdCreat int
	ChangeTime   int64
	lock         *sync.Mutex
	wg           *sync.WaitGroup
}

func NewDailyAd() *DailyAd {
	return &DailyAd{
		DailyAdCreat: 0,
		ChangeTime:   0,
		lock:         &sync.Mutex{},
		wg:           &sync.WaitGroup{},
	}
}

func (d *DailyAd) BeforeAPIEvent(adCtx *types.AdControllerCtx, srv service.AdService) error {
	ctx := adCtx.Ctx
	err := d.event(ctx, srv)
	if err != nil {
		return err
	}
	if d.DailyAdCreat >= 3000 {
		return errors.New("daily limit of 3000 data entries")
	}
	return nil
}

func (d *DailyAd) AfterAPIEvent(adCtx *types.AdControllerCtx, srv service.AdService) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.DailyAdCreat++
	return nil
}

func (d *DailyAd) Pipeline() mongo.Pipeline {
	// Get the start and end timestamps of the current day's ads
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	return mongo.Pipeline{
		//nolint:govet
		{{"$match", bson.D{
			{"timestamp", bson.D{
				{"$gte", startOfDay.Unix()}, // Match documents where timestamp is greater than or equal to start of day
				{"$lt", endOfDay.Unix()},    // and less than end of day
			}},
		}}},
		//nolint:composites
		{{"$group", bson.D{
			{"_id", nil},
			{"daily_ad", bson.D{
				{"$count", bson.D{}},
			}},
		}}},
	}
}

func (d *DailyAd) InitEvent(adCtx *types.AdControllerCtx, srv service.AdService) error {
	ctx := adCtx.Ctx
	err := d.event(ctx, srv)
	if err != nil {
		return err
	}
	return nil
}

func (d *DailyAd) event(ctx context.Context, srv service.AdService) error {
	var result struct {
		DailyAd int `bson:"daily_ad"`
	}

	err := srv.Aggregate(context.TODO(), d, &result)
	if err != nil {
		return err
	}
	d.DailyAdCreat = result.DailyAd
	log.Println("today creat ad count", d.DailyAdCreat)
	return nil
}

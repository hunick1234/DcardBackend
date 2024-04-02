package controller

import (
	"context"
	"errors"
	"log"
	"net/http"
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
	temp         int
	dailyLimit   int
	lock         *sync.Mutex
}

func NewDailyAd() *DailyAd {
	return &DailyAd{
		DailyAdCreat: 0,
		temp:         0,
		ChangeTime:   0,
		dailyLimit:   10000,
		lock:         &sync.Mutex{},
	}
}

func (d *DailyAd) BeforeAPIEvent(adCtx *types.AdControllerCtx, srv service.AdService) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	log.Println("daily ad: ", d.DailyAdCreat)
	if d.DailyAdCreat >= d.dailyLimit {
		log.Println("daily limit of 3000 data entries")
		adCtx.W.StausCode = http.StatusTooManyRequests
		adCtx.Err = errors.New("daily limit of 3000 data entries")
		return nil
	}

	//wait the
	for d.temp >= d.dailyLimit && d.DailyAdCreat < d.dailyLimit {
		time.Sleep(10 * time.Millisecond)
	}
	d.temp++
	return nil
}

func (d *DailyAd) AfterAPIEvent(adCtx *types.AdControllerCtx, srv service.AdService) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	if adCtx.W.StausCode < 200 || adCtx.W.StausCode >= 300 {
		d.temp--
		return nil
	}

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
	d.startDailyReset(srv)
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

func (d *DailyAd) startDailyReset(srv service.AdService) {
	ticker := time.NewTicker(time.Until(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, 0, 0, 0, 0, time.Now().Location())))
	go func() {
		for range ticker.C {
			d.lock.Lock()
			d.DailyAdCreat = 0
			d.temp = 0
			d.lock.Unlock()
			d.event(context.Background(), srv)
		}
	}()
}

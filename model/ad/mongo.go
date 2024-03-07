package ad

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/hunick1234/DcardBackend/config"
	"github.com/hunick1234/DcardBackend/model"
	"github.com/hunick1234/DcardBackend/storage"
)

const collectionName string = "ad"
const databaseName string = "dcard"

var AdService *adService
var syncOnce sync.Once

type adService struct {
	dbClient storage.Storager
}

func init() {
	// implementation of init method
	var _ model.Storager[AD, AdQuery] = (*adService)(nil)
	AdService = DeafultAdService()
}

func NewAdService(dbClient storage.Storager) *adService {
	return &adService{
		dbClient: dbClient,
	}
}

func DeafultAdService() *adService {
	syncOnce.Do(func() {
		deafult := &config.MongoCfg{
			URI: "mongodb://localhost:27017",
			DB:  databaseName,
		}
		storager, err := storage.NewMongoConn(deafult)
		if err != nil {
			log.Fatal(err)
		}
		AdService = NewAdService(storager)
	})

	return AdService
}

// findByFilter implements model.Storager.
func (service *adService) FindByFilter(ctx context.Context, adQuery *AdQuery) (*[]AD, error) {

	if service.dbClient == nil {
		return nil, fmt.Errorf("check you DB connection, it's nil")
	}
	start := time.Now()
	collection, err := service.dbClient.GetCollection(collectionName)
	if err != nil {
		return nil, err
	}

	filter := adQuery.Pipeline()

	// 執行查詢
	cur, err := collection.Aggregate(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	var results []AD
	// Iterate through the cursor allowing to decode documents
	var t storage.Translate
	t.With(&AD{}).Decodes(cur, &results)

	elapsed := time.Since(start)
	fmt.Printf("Search time: %s\n", elapsed)

	return &results, nil
}

// Store implements model.Storager.
func (service *adService) Store(ad *AD) error {
	ad.Timestamp = time.Now().Unix()

	collection, err := service.dbClient.GetCollection(collectionName)
	if err != nil {
		return err
	}
	_, err = collection.InsertOne(context.Background(), ad)
	if err != nil {
		return err
	}
	return nil
}

func (service *adService) Getlivead() (int, error) {
	return 0, nil
}

func (service *adService) GetDailyCreatAd() (int, error) {
	return 0, nil
}

func (service *adService) Aggregate(ctx context.Context, filter model.Filter, result any) error {
	pipe := filter.Pipeline()
	collection, err := service.dbClient.GetCollection(collectionName)
	if err != nil {
		return err
	}

	cur, err := collection.Aggregate(ctx, pipe)
	if err != nil {
		return err
	}

	var t storage.Translate
	if reflect.ValueOf(result).Kind() == reflect.Slice {
		t.With(result).Decodes(cur, result)
	} else {
		t.With(result).Decode(cur, result)
	}

	return nil
}

package ad

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hunick1234/DcardBackend/model"
	"github.com/hunick1234/DcardBackend/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var AdService *adService
var collectionName = "ad"
var databaseName = "dcard"

type adService struct {
	DbClient *storage.MongoDB
}

func init() {
	// implementation of init method
	var _ model.Storage[AD, AdQuery] = (*adService)(nil)
	AdService = DeafultAdService()
}

func NewAdService(dbClient *storage.MongoDB) *adService {
	return &adService{
		DbClient: dbClient,
	}
}

func DeafultAdService() *adService {
	dbClient, err := storage.Connect(options.Client().ApplyURI(storage.MongoAddress), databaseName)
	dbClient.CollectionName = collectionName
	if err != nil {
		log.Fatal(err)
	}

	return &adService{
		DbClient: dbClient,
	}
}

// findByFilter implements model.Storage.
func (service *adService) FindByFilter(ctx context.Context, adQuery AdQuery) ([]*AD, error) {
	if service.DbClient == nil {
		return nil, fmt.Errorf("check you DB connection, it's nil")
	}
	start := time.Now()
	fmt.Println("Connected to MongoDB!")
	collection := service.DbClient.DB.Collection(service.DbClient.CollectionName)
	// ESR
	var filter bson.D = bson.D{}

	if len(adQuery.Gender) > 0 {
		filter = append(filter, bson.D{
			{"$or", bson.A{
				bson.D{{"conditios.gender", []string{}}},
				bson.D{{"conditions.gender", adQuery.Gender}},
			}},
		}...)
	}

	if len(adQuery.Platform) > 0 {
		filter = append(filter, bson.D{
			{"$or", bson.A{
				bson.D{{"conditios.platform", []string{}}},
				bson.D{{"conditions.platform", adQuery.Platform}},
			}},
		}...)
	}
	if len(adQuery.Country) > 0 {
		filter = append(filter, bson.D{
			{"$or", bson.A{
				bson.D{{"conditios.country", []string{}}},
				bson.D{{"conditions.country", adQuery.Country}},
			}},
		}...)
	}

	sort := bson.D{{"end_at", -1}}

	if adQuery.Age > 0 {
		rangeFilter := bson.D{
			{"$and", bson.A{
				bson.D{{"conditions.age_start", bson.D{{"$lte", adQuery.Age}}}},
				bson.D{{"conditions.age_end", bson.D{{"$gte", adQuery.Age}}}},
			}},
		}
		filter = append(filter, rangeFilter...)
	}
	fmt.Println(filter)

	// 執行查詢
	cur, err := collection.Find(context.TODO(), filter, options.Find().SetSort(sort))
	if err != nil {
		log.Fatal(err)
	}

	var resultAds []*AD
	for cur.Next(context.TODO()) {
		var ad AD
		err := cur.Decode(&ad)
		if err != nil {
			log.Fatal(err)
		}
		resultAds = append(resultAds, &ad)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(context.TODO())

	elapsed := time.Since(start)
	fmt.Printf("Search time: %s\n", elapsed)

	return resultAds, nil
}

// Store implements model.Storage.
func (a *adService) Store(ad *AD) error {
	ad.Timestamp = time.Now().Unix()

	_, err := a.DbClient.DB.Collection(a.DbClient.CollectionName).InsertOne(context.Background(), ad)
	if err != nil {
		return err
	}
	return nil
}

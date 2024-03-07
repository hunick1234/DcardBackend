package storage

import (
	"context"
	"log"
	"reflect"

	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
}

type Translate struct {
	reflectType reflect.Type
}

func (t *Translate) Decode(cur *mongo.Cursor, result interface{}) {
	// Check if result is of type *interface{}
	val := reflect.ValueOf(result)
	if val.Kind() != reflect.Ptr {
		log.Fatal("result must be a pointer")
	}

	if cur.Next(context.Background()) {
		err := cur.Decode(result)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("No documents found")
	}
}
func (t *Translate) Decodes(cur *mongo.Cursor, results any) {
	// Check if results is a pointer to a slice
	val := reflect.ValueOf(results)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
		log.Fatal("results must be a pointer to a slice")
	}
	// Prepare slice of the specific type defined by t.elemType
	resultsSlice := reflect.MakeSlice(reflect.SliceOf(t.reflectType), 0, 0)
	// Iterate over the cursor
	for cur.Next(context.Background()) {
		// Create a new element of the type to decode into
		elemPtr := reflect.New(t.reflectType) // Note: ensure t.elemType is not a pointer itself
		err := cur.Decode(elemPtr.Interface())
		if err != nil {
			log.Fatal(err)
		}
		// Append decoded element to the results slice
		resultsSlice = reflect.Append(resultsSlice, elemPtr.Elem())
	}

	// Set the decoded slice to the results using reflection
	val.Elem().Set(resultsSlice)
}

func (t *Translate) With(value interface{}) *Translate {
	reflectType := reflect.TypeOf(value)
	if reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}
	t.reflectType = reflectType
	return t
}

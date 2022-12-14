package dbHelper

import (
	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	conf "goex/config"
	db "goex/db/mongoPools"
)

var ErrorInsertFailed = errors.New("insert failed with no error")

type Recorder interface {
	GetCollectionName() string
	GetID() interface{}
	SetID(id interface{})
	SetIsDocumented(bool)
	GetIsDocumented() bool
}
type Transaction struct {
	WMongo     *db.MongoWriteDB
	Connection *mongo.Database
}
type DecoderMap func(m map[string]interface{}) error

func StartTransaction() *Transaction {
	tr := new(Transaction)
	tr.WMongo = db.GetWriteDB()
	tr.Connection = tr.WMongo.GetConnection()
	return tr
}

func (t *Transaction) EndTransaction(f func(sessionContext mongo.SessionContext) (result interface{}, err error)) (interface{}, error) {

	wc := writeconcern.New(writeconcern.W(1))
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	session, err := t.Connection.Client().StartSession()
	if err != nil {
		return nil, err
	}
	resp, err := session.WithTransaction(context.Background(), f, txnOpts)
	if err != nil {
		//err=session.AbortTransaction(context.Background())
	}

	defer session.EndSession(context.Background())
	defer t.WMongo.Release(t.Connection)
	return resp, err

}

func GetCollection(collectionName string) (*mongo.Collection, func()) {
	mongoDB := db.GetWriteDB()
	conn := mongoDB.GetConnection()
	return conn.Collection(collectionName), func() { mongoDB.Release(conn) }
}

func Insert(ctx context.Context, record Recorder, res chan error) {

	var collection *mongo.Collection
	if cs, ok := ctx.(mongo.SessionContext); ok {
		collection = cs.Client().Database(conf.GetMongodbName()).Collection(record.GetCollectionName())
	} else {
		var release func()
		collection, release = GetCollection(record.GetCollectionName())
		defer release()
	}

	insertResult, err := collection.InsertOne(ctx, record)
	if err != nil {
		res <- err
		return
	}
	if insertResult == nil {
		res <- ErrorInsertFailed
		return
	}
	record.SetID(insertResult.InsertedID)
	record.SetIsDocumented(true)
	res <- nil
}
func UpdateGo(ctx context.Context, record Recorder) chan error {
	res := make(chan error)
	go Update(ctx, record, res)
	return res

}
func UpdateMapGo(ctx context.Context, collectionName string, id *primitive.ObjectID, newData map[string]interface{}) chan error {
	res := make(chan error)
	go func() {
		var collection *mongo.Collection
		if cs, ok := ctx.(mongo.SessionContext); ok {
			collection = cs.Client().Database(conf.GetMongodbName()).Collection(collectionName)
		} else {
			var release func()
			collection, release = GetCollection(collectionName)
			defer release()
		}

		updateFilter := bson.D{{Key: "$set", Value: newData}}
		updatedResult, err := collection.UpdateOne(ctx, bson.D{{"_id", id}}, updateFilter)
		if err != nil {
			res <- err
			return
		}
		if updatedResult == nil || updatedResult.MatchedCount == 0 {
			res <- mongo.ErrNoDocuments
			return
		}
		res <- nil
	}()
	return res
}
func Update(ctx context.Context, record Recorder, res chan error) {

	var collection *mongo.Collection
	if cs, ok := ctx.(mongo.SessionContext); ok {
		collection = cs.Client().Database(conf.GetMongodbName()).Collection(record.GetCollectionName())
	} else {
		var release func()
		collection, release = GetCollection(record.GetCollectionName())
		defer release()
	}

	updateFilter := bson.D{{Key: "$set", Value: record}}
	updatedResult, err := collection.UpdateOne(ctx, bson.D{{"_id", record.GetID()}}, updateFilter)
	if err != nil {
		res <- err
		return
	}
	if updatedResult == nil || updatedResult.MatchedCount == 0 {
		res <- mongo.ErrNoDocuments
		return
	}
	record.SetIsDocumented(true)
	res <- nil
}
func UpdateMany(ctx context.Context, collectionName string, filter interface{}, update interface{}, MatchCount *int64, options ...*options.UpdateOptions) chan error {
	res := make(chan error)
	go func() {
		var collection *mongo.Collection
		if cs, ok := ctx.(mongo.SessionContext); ok {
			collection = cs.Client().Database(conf.GetMongodbName()).Collection(collectionName)
		} else {
			var release func()
			collection, release = GetCollection(collectionName)
			defer release()
		}
		updatedResult, err := collection.UpdateMany(ctx, filter, update, options...)
		if err != nil {
			res <- err
		}
		if updatedResult == nil {
			res <- mongo.ErrNoDocuments
		} else {
			*MatchCount = updatedResult.MatchedCount
		}
		close(res)
	}()

	return res
}
func SaveGo(ctx context.Context, record Recorder) chan error {
	res := make(chan error)
	Save(ctx, record, res)
	return res
}

func SaveSync(ctx context.Context, record Recorder) (err error) {
	res := make(chan error)
	Save(ctx, record, res)
	return <-res
}

func Save(ctx context.Context, record Recorder, res chan error) {
	if record.GetIsDocumented() {
		go Update(ctx, record, res)
	} else {
		go Insert(ctx, record, res)
	}
}
func FindOneMap(ctx context.Context, query *bson.D, recorder map[string]interface{}, collectionName string, res chan error, options ...*options.FindOneOptions) {
	var collection *mongo.Collection
	if cs, ok := ctx.(mongo.SessionContext); ok {
		collection = cs.Client().Database(conf.GetMongodbName()).Collection(collectionName)
	} else {
		var release func()
		collection, release = GetCollection(collectionName)
		defer release()
	}
	one := collection.FindOne(ctx, query, options...)
	err := one.Decode(recorder)
	res <- err
}
func FindOne(ctx context.Context, query *bson.D, record Recorder, res chan error, options ...*options.FindOneOptions) {
	var collection *mongo.Collection
	if cs, ok := ctx.(mongo.SessionContext); ok {
		collection = cs.Client().Database(conf.GetMongodbName()).Collection(record.GetCollectionName())
	} else {
		var release func()
		collection, release = GetCollection(record.GetCollectionName())
		defer release()
	}
	one := collection.FindOne(ctx, query, options...)
	err := one.Decode(record)
	if err == nil {
		record.SetIsDocumented(true)
	}
	res <- err
}

func FindOneGo(ctx context.Context, query *bson.D, record Recorder, options ...*options.FindOneOptions) chan error {
	res := make(chan error)
	go FindOne(ctx, query, record, res, options...)
	return res
}

func FindOneMapGo(ctx context.Context, query *bson.D, record map[string]interface{}, collectionName string, options ...*options.FindOneOptions) chan error {
	res := make(chan error)
	go FindOneMap(ctx, query, record, collectionName, res, options...)
	return res
}

func FindOneSync(ctx context.Context, query *bson.D, record Recorder) (err error) {
	res := make(chan error)
	go FindOne(ctx, query, record, res)
	go FindOne(ctx, query, record, res)
	return <-res
}

type Decoder func(model Recorder) error

func FindAllOnMap(ctx context.Context, collectionName string, query *bson.D, result chan DecoderMap, opts ...*options.FindOptions) {
	var collection *mongo.Collection
	if cs, ok := ctx.(mongo.SessionContext); ok {
		collection = cs.Client().Database(conf.GetMongodbName()).Collection(collectionName)
	} else {
		var release func()
		collection, release = GetCollection(collectionName)
		defer release()
	}
	cur, err := collection.Find(ctx, query, opts...)
	if err != nil {
		result <- func(_ map[string]interface{}) error { return err }
		close(result)
		return
	}

	defer func() { _ = cur.Close(ctx) }()
	for cur.Next(ctx) {
		result <- func(cursor mongo.Cursor) DecoderMap {
			return func(model map[string]interface{}) error {
				err := cursor.Decode(model)
				if err != nil {
					return err
				}
				return nil
			}
		}(*cur)
	}
	close(result)
}
func FindAll(ctx context.Context, collectionName string, query *bson.D, result chan Decoder, opts ...*options.FindOptions) {
	var collection *mongo.Collection
	if cs, ok := ctx.(mongo.SessionContext); ok {
		collection = cs.Client().Database(conf.GetMongodbName()).Collection(collectionName)
	} else {
		var release func()
		collection, release = GetCollection(collectionName)
		defer release()
	}
	cur, err := collection.Find(ctx, query, opts...)
	if err != nil {
		result <- func(_ Recorder) error { return err }
		close(result)
		return
	}

	defer func() { _ = cur.Close(ctx) }()
	for cur.Next(ctx) {
		result <- func(cursor mongo.Cursor) Decoder {
			return func(model Recorder) error {
				err := cursor.Decode(model)
				if err != nil {
					return err
				}
				model.SetIsDocumented(true)
				return nil
			}
		}(*cur)
	}
	close(result)
}

func FindAllGo(ctx context.Context, collectionName string, query *bson.D, opts ...*options.FindOptions) chan Decoder {
	res := make(chan Decoder)
	go FindAll(ctx, collectionName, query, res, opts...)
	return res
}
func FindAllOnMapGo(ctx context.Context, collectionName string, query *bson.D, opts ...*options.FindOptions) chan DecoderMap {
	res := make(chan DecoderMap)
	go FindAllOnMap(ctx, collectionName, query, res, opts...)
	return res
}

func FindByID(ctx context.Context, record Recorder, id interface{}, res chan error) {
	var collection *mongo.Collection
	if cs, ok := ctx.(mongo.SessionContext); ok {
		collection = cs.Client().Database(conf.GetMongodbName()).Collection(record.GetCollectionName())
	} else {
		var release func()
		collection, release = GetCollection(record.GetCollectionName())
		defer release()
	}
	one := collection.FindOne(ctx, &bson.D{{Key: "_id", Value: id}})
	err := one.Err()
	if err != nil {
		res <- err
		return
	}
	err = one.Decode(record)
	record.SetIsDocumented(true)
	res <- err
}

func FindByIDGo(ctx context.Context, record Recorder, id interface{}) chan error {
	res := make(chan error)
	go FindByID(ctx, record, id, res)
	return res
}

func FindByIDSync(ctx context.Context, record Recorder, id interface{}) error {
	res := make(chan error)
	go FindByID(ctx, record, id, res)
	return <-res
}

func DeleteOneGo(ctx context.Context, collectionName string, b *bson.D, opts ...*options.DeleteOptions) chan error {
	res := make(chan error)
	go func() {
		var collection *mongo.Collection
		if cs, ok := ctx.(mongo.SessionContext); ok {
			collection = cs.Client().Database(conf.GetMongodbName()).Collection(collectionName)
		} else {
			var release func()
			collection, release = GetCollection(collectionName)
			defer release()
		}

		_, err := collection.DeleteOne(ctx, b, opts...)
		if err != nil {
			res <- err
		}
		close(res)
	}()
	return res
}

func DeleteRecordGo(ctx context.Context, record Recorder) chan error {
	return DeleteOneGo(ctx, record.GetCollectionName(), &bson.D{{Key: "_id", Value: record.GetID()}})
}

func DeleteManyRecordGo(ctx context.Context, records []Recorder, deletedCount *int64) chan error {
	query := bson.A{}
	var collectionName string
	c := 0
	for _, r := range records {
		if c == 0 {
			collectionName = r.GetCollectionName()
		}
		if r.GetCollectionName() != collectionName {
			e := make(chan error)
			go func() {
				e <- errors.New("records from different collection not supported")
			}()
			return e
		}
		collectionName = r.GetCollectionName()
		query = append(query, bson.E{Key: "_id", Value: r.GetID()})
		c++
	}

	return DeleteManyGo(ctx, collectionName, &bson.D{{Key: "$or", Value: query}}, deletedCount)
}

func DeleteManyGo(ctx context.Context, collectionName string, b *bson.D, deletedCount *int64, opts ...*options.DeleteOptions) chan error {
	result := make(chan error)
	go func() {
		var collection *mongo.Collection
		if cs, ok := ctx.(mongo.SessionContext); ok {
			collection = cs.Client().Database(conf.GetMongodbName()).Collection(collectionName)
		} else {
			var release func()
			collection, release = GetCollection(collectionName)
			defer release()
		}

		delCount, err := collection.DeleteMany(ctx, b, opts...)
		if err != nil {
			result <- err
		} else {
			if deletedCount != nil {
				*deletedCount = delCount.DeletedCount
			}

		}
		close(result)
	}()
	return result
}

type AdvancedQuery struct {
	_collection string
	_pipeline   []bson.M
	_limit      bson.M
	_skip       bson.M
	_c          context.Context
}

// func NewAdvancedQuery(collection string, fields string, query string, pagination *Pagination, project *bson.M) *AdvancedQuery {
// 	option := pagination.CreateFindOption()
// 	aq := &AdvancedQuery{_collection: collection}
// 	queryParams := strings.Split(query, " ")
// 	fieldArray := strings.Split(fields, " ")
// 	signedFieldsArray := make([]bson.M, len(fieldArray))
// 	for i, f := range fieldArray {
// 		signedFieldsArray[i] = bson.M{"$toString": "$" + f}
// 	}

// 	d := make([]interface{}, 2)
// 	d[0] = bson.M{"$toLower": "$" + fieldArray[0]}
// 	d[1] = queryParams[0]

// 	addFields := bson.M{
// 		"_aq_selectedFields": bson.M{"$concat": signedFieldsArray},
// 		//"_aq_score":          bson.M{"$indexOfCP": d},
// 	}

// 	matchConditions := make([]bson.M, len(queryParams))
// 	for i, q := range queryParams {
// 		matchConditions[i] = bson.M{"_aq_selectedFields": primitive.Regex{
// 			Pattern: "" + q + "",
// 			Options: "gi",
// 		}}
// 	}

// 	match := bson.M{"$and": matchConditions}
// 	aq._pipeline = []bson.M{
// 		{"$addFields": addFields},
// 		{"$match": match},
// 		//{"$sort": bson.M{"_aq_score": 1}},
// 	}
// 	if project != nil {
// 		aq._pipeline = append(aq._pipeline, *project)
// 	}
// 	aq._limit = bson.M{"$limit": option.Limit}
// 	aq._skip = bson.M{"$skip": option.Skip}
// 	return aq
// }

func (aq *AdvancedQuery) QueryGo(ctx context.Context) chan Decoder {
	res := make(chan Decoder)
	go func() {
		cur, err := aq.Query(ctx)
		if err != nil {
			res <- func(_ Recorder) error {
				return err
			}
			close(res)
			return
		}
		for cur.Next(ctx) {
			res <- func(cursor mongo.Cursor) Decoder {
				return func(model Recorder) error {
					err := cursor.Decode(model)
					if err != nil {
						println(err.Error())
						return err
					}
					model.SetIsDocumented(true)
					return nil
				}
			}(*cur)
		}
		close(res)
	}()
	return res
}
func (aq *AdvancedQuery) CountGo(ctx context.Context, count *int64) chan error {
	res := make(chan error)
	go func() {
		countRes, err := aq.Count(ctx)
		if err != nil {
			res <- err
		} else {
			*count = (countRes)
			res <- nil
		}
	}()
	return res
}
func (aq *AdvancedQuery) Query(ctx context.Context) (*mongo.Cursor, error) {
	var collection *mongo.Collection
	if cs, ok := ctx.(mongo.SessionContext); ok {
		collection = cs.Client().Database(conf.GetMongodbName()).Collection(aq._collection)
	} else {
		var release func()
		collection, release = GetCollection(aq._collection)
		defer release()
	}
	pipe := aq._pipeline
	pipe = append(pipe, aq._limit)
	pipe = append(pipe, aq._skip)
	return collection.Aggregate(ctx, pipe)
}
func (aq *AdvancedQuery) Count(ctx context.Context) (int64, error) {
	var collection *mongo.Collection
	if cs, ok := ctx.(mongo.SessionContext); ok {
		collection = cs.Client().Database(conf.GetMongodbName()).Collection(aq._collection)
	} else {
		var release func()
		collection, release = GetCollection(aq._collection)
		defer release()
	}
	pipe := aq._pipeline
	pipe = append(pipe, bson.M{"$count": "_aq_count"})

	return collection.CountDocuments(ctx, pipe)
}

func AdvanceQueryCursor(ctx context.Context, collectionName string, fields string, query string, paginationOptions ...int) (*mongo.Cursor, error) {
	var collection *mongo.Collection
	if cs, ok := ctx.(mongo.SessionContext); ok {
		collection = cs.Client().Database(conf.GetMongodbName()).Collection(collectionName)
	} else {
		var release func()
		collection, release = GetCollection(collectionName)
		defer release()
	}
	queryParams := strings.Split(query, " ")
	fieldArray := strings.Split(fields, " ")
	signedFieldsArray := make([]bson.M, len(fieldArray))
	for i, f := range fieldArray {
		signedFieldsArray[i] = bson.M{"$toString": "$" + f}
	}

	d := make([]interface{}, 2)
	d[0] = bson.M{"$toLower": "$" + fieldArray[0]}
	d[1] = queryParams[0]

	addFields := bson.M{
		"_aq_selectedFields": bson.M{"$concat": signedFieldsArray},
		"_aq_score":          bson.M{"$indexOfCP": d},
	}

	matchConditions := make([]bson.M, len(queryParams))
	for i, q := range queryParams {
		matchConditions[i] = bson.M{"_aq_selectedFields": primitive.Regex{
			Pattern: "" + q + "",
			Options: "gi",
		}}
	}

	match := bson.M{"$and": matchConditions}

	limit := 100
	offset := 0
	if len(paginationOptions) > 0 {
		limit = paginationOptions[0]
	}
	if len(paginationOptions) > 1 {
		offset = paginationOptions[1]
	}

	return collection.Aggregate(ctx, []bson.M{
		{"$addFields": addFields},
		{"$match": match},
		{"$sort": bson.M{"_aq_score": 1}},
		{"$limit": limit},
		{"$skip": offset},
	})

}

func AdvanceQueryG(ctx context.Context, collectionName string, fields string, query string, paginationOptions ...int) chan Decoder {
	res := make(chan Decoder)
	go func() {
		cur, err := AdvanceQueryCursor(ctx, collectionName, fields, query, paginationOptions...)
		if err != nil {
			res <- func(_ Recorder) error {
				return err
			}
			close(res)
			return
		}
		for cur.Next(ctx) {
			res <- func(cursor mongo.Cursor) Decoder {
				return func(model Recorder) error {
					err := cursor.Decode(model)
					if err != nil {
						return err
					}
					model.SetIsDocumented(true)
					return nil
				}
			}(*cur)
		}
		close(res)
	}()
	return res
}

func Count(ctx context.Context, CollectionName string, query interface{}, count *int64, res chan error, opts ...*options.CountOptions) {
	var collection *mongo.Collection
	if cs, ok := ctx.(mongo.SessionContext); ok {
		collection = cs.Client().Database(conf.GetMongodbName()).Collection(CollectionName)
	} else {
		var release func()
		collection, release = GetCollection(CollectionName)
		defer release()
	}
	documents, err := collection.CountDocuments(ctx, query, opts...)
	if err != nil {
		res <- err
		return
	}
	*count = documents
	res <- nil
}

func CountGo(c context.Context, CollectionName string, query interface{}, count *int64, opts ...*options.CountOptions) chan error {
	res := make(chan error)
	go Count(c, CollectionName, query, count, res, opts...)
	return res
}
func CountCollectionDocGo(ctx context.Context, CollectionName string, query interface{}, count *int64, opts ...*options.CountOptions) chan error {
	res := make(chan error)
	go func() {
		mongoDB := db.GetWriteDB()
		dbase := mongoDB.GetConnection()
		result := dbase.RunCommand(ctx, bson.M{"collStats": CollectionName})
		var document bson.M
		err := result.Decode(&document)
		if err != nil {
			res <- err
			return
		}
		*count = int64(document["count"].(int32))
		res <- nil

	}()
	//go Count(c, CollectionName, query, count, res, opts...)
	return res
}
func CountSync(c context.Context, CollectionName string, query interface{}, count *int64, opts ...*options.CountOptions) error {
	res := make(chan error)
	go Count(c, CollectionName, query, count, res, opts...)
	return <-res
}

func AggregateGo(ctx context.Context, collectionName string, pipe mongo.Pipeline, opts ...*options.AggregateOptions) chan Decoder {
	result := make(chan Decoder)
	go func() {

		var collection *mongo.Collection
		if cs, ok := ctx.(mongo.SessionContext); ok {
			collection = cs.Client().Database(conf.GetMongodbName()).Collection(collectionName)
		} else {
			var release func()
			collection, release = GetCollection(collectionName)
			defer release()
		}
		cur, err := collection.Aggregate(ctx, pipe, opts...)
		if err != nil {
			result <- func(_ Recorder) error { return err }
			close(result)
			return
		}

		defer func() { _ = cur.Close(ctx) }()
		for cur.Next(ctx) {
			result <- func(cursor mongo.Cursor) Decoder {
				return func(model Recorder) error {
					err := cursor.Decode(model)
					if err != nil {
						return err
					}
					model.SetIsDocumented(true)
					return nil
				}
			}(*cur)
		}
		close(result)

	}()
	return result
}
func AggregateMapGo(ctx context.Context, collectionName string, pipe mongo.Pipeline, opts ...*options.AggregateOptions) chan DecoderMap {
	result := make(chan DecoderMap)
	go func() {

		var collection *mongo.Collection
		if cs, ok := ctx.(mongo.SessionContext); ok {
			collection = cs.Client().Database(conf.GetMongodbName()).Collection(collectionName)
		} else {
			var release func()
			collection, release = GetCollection(collectionName)
			defer release()
		}
		cur, err := collection.Aggregate(ctx, pipe, opts...)
		if err != nil {
			result <- func(_ map[string]interface{}) error { return err }
			close(result)
			return
		}

		defer func() { _ = cur.Close(ctx) }()
		for cur.Next(ctx) {
			result <- func(cursor mongo.Cursor) DecoderMap {
				return func(model map[string]interface{}) error {
					err := cursor.Decode(model)
					if err != nil {
						return err
					}
					return nil
				}
			}(*cur)
		}
		close(result)

	}()
	return result
}

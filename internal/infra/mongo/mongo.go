package mongo

import (
	"context"

	"github.com/skyrocketOoO/RBAC-server/domain"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	db         = viper.GetString("mongo.db")
	collection = viper.GetString("mongo.collection")
)

type MongoRepository struct {
	client *mongo.Client
}

func NewMongoRepository(client *mongo.Client) (*MongoRepository, error) {
	return &MongoRepository{
		client: client,
	}, nil
}

func (r *MongoRepository) Ping(c context.Context) error {
	return r.client.Ping(c, readpref.Primary())
}

func (r *MongoRepository) Get(c context.Context, filter domain.Edge, queryMode bool) (
	[]domain.Edge, error) {
	col := r.client.Database(db).Collection(collection)
	edges := []domain.Edge{}
	if queryMode {
		if filter == (domain.Edge{}) {
			// get all records
			cursor, err := col.Find(c, bson.D{})
			if err != nil {
				return nil, err
			}
			defer cursor.Close(c)
			if err := cursor.All(c, &edges); err != nil {
				return nil, err
			}
		} else {
			cursor, err := col.Find(c, rmZeroVal(filter))
			if err != nil {
				return nil, err
			}
			defer cursor.Close(c)
			if err := cursor.All(c, &edges); err != nil {
				return nil, err
			}
		}
	} else {
		cursor, err := col.Find(c, filter)
		if err != nil {
			return nil, err
		}
		defer cursor.Close(c)
		if err := cursor.All(c, &edges); err != nil {
			return nil, err
		}
		if len(edges) == 0 {
			return nil, domain.ErrRecordNotFound
		} else if len(edges) > 1 {
			return nil, domain.ErrDuplicateRecord
		}
	}
	return edges, nil
}

func (r *MongoRepository) Create(c context.Context, edge domain.Edge) error {
	col := r.client.Database(db).Collection(collection)
	_, err := col.InsertOne(c, edge)
	return err
}

func (r *MongoRepository) Delete(c context.Context, edge domain.Edge,
	queryMode bool) error {
	col := r.client.Database(db).Collection(collection)
	if queryMode {
		_, err := col.DeleteMany(c, rmZeroVal(edge))
		return err
	} else {
		if _, err := r.Get(c, edge, false); err != nil {
			return err
		}
		_, err := col.DeleteOne(c, edge)
		return err
	}
}

func (r *MongoRepository) ClearAll(c context.Context) error {
	col := r.client.Database(db).Collection(collection)
	_, err := col.DeleteMany(c, bson.M{})
	return err
}

func rmZeroVal(filter domain.Edge) bson.M {
	m := bson.M{}
	if filter.VNs != "" {
		m["v_ns"] = filter.VNs
	}
	if filter.VName != "" {
		m["v_name"] = filter.VName
	}
	if filter.Rel != "" {
		m["rel"] = filter.Rel
	}
	if filter.UNs != "" {
		m["u_ns"] = filter.UNs
	}
	if filter.UName != "" {
		m["u_name"] = filter.UName
	}
	return m
}

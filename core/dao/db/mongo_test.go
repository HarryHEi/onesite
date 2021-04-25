package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"

	"onesite/core/config"
)

func TestInitMongo(t *testing.T) {
	config.CfgRootPath = "../../../configs"

	mg, err := NewMongo()
	require.Nil(t, err)
	collection := mg.Db.Collection("test.demo")
	inserted, err := collection.InsertOne(context.Background(), bson.D{{"name", "test"}})
	require.Nil(t, err)
	queryRes := collection.FindOne(context.Background(), bson.D{})
	require.Nil(t, queryRes.Err())
	deleteRes, err := collection.DeleteOne(context.Background(), bson.D{{"_id", inserted.InsertedID}})
	require.Nil(t, err)
	require.NotZero(t, deleteRes.DeletedCount)
	err = collection.Drop(context.Background())
	require.Nil(t, err)
}

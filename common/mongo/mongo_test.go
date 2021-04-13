package mongo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"

	mongo2 "onesite/common/mongo"
)

func TestInitMongo(t *testing.T) {
	_, err := mongo2.GetMongo()
	require.NotNil(t, err)
	err = mongo2.InitMongo()
	require.Nil(t, err)
	mg, err := mongo2.GetMongo()
	require.Nil(t, err)
	collection := mg.Collection("test.demo")
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

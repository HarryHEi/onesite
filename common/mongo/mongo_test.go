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
	client, err := mongo2.GetMongo()
	require.Nil(t, err)
	database := client.Database("onesite")
	collection := database.Collection("chat.message")
	res := collection.FindOne(context.Background(), bson.D{})
	require.Nil(t, res.Err())
}

package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"onesite/core/model"
)

func (dao *Dao) CreateArticle(article *model.Article) error {
	collection := dao.Mongo.Db.Collection(model.ArticleCollectionName)
	article.Creation = time.Now()
	_, err := collection.InsertOne(context.Background(), article)
	return err
}

func (dao *Dao) QueryArticleView(page, pageSize int) (int64, []model.Article, error) {
	if page <= 0 {
		page = 1
	}

	limit := int64(pageSize)
	offset := int64((page - 1) * pageSize)

	collection := dao.Mongo.Db.Collection(model.ArticleCollectionName)
	opts := options.FindOptions{
		Limit:      &limit,
		Skip:       &offset,
		Sort:       bson.D{{"creation", -1}},
		Projection: bson.D{{"title", 1}, {"creation", 1}},
	}
	count, err := collection.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		return 0, nil, err
	}
	cur, err := collection.Find(context.Background(), bson.D{}, &opts)
	if err != nil {
		return 0, nil, err
	}
	var articles []model.Article
	err = cur.All(context.Background(), &articles)
	if err != nil {
		return 0, nil, err
	}
	return count, articles, nil
}

func (dao *Dao) QueryArticleDetail(pk string) (*model.Article, error) {
	collection := dao.Mongo.Db.Collection(model.ArticleCollectionName)
	objectId, err := primitive.ObjectIDFromHex(pk)
	if err != nil {
		return nil, err
	}
	res := collection.FindOne(context.Background(), bson.D{{"_id", objectId}})
	if res.Err() != nil {
		return nil, res.Err()
	}
	var article model.Article
	err = res.Decode(&article)
	if err != nil {
		return nil, err
	}
	return &article, nil
}

package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ArticleCollectionName = "blog.article"
)

type Article struct {
	ID primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	// Author is user's username
	Author   string    `json:"author" bson:"author"`
	Title    string    `json:"title" bson:"title"`
	Document string    `json:"document" bson:"document"`
	Comments []Comment `json:"comments" bson:"comments"`
	Creation time.Time `json:"creation" bson:"creation,omitempty"`
}

type Comment struct {
	Author     string    `json:"author" bson:"author"`
	AuthorName string    `json:"author_name" bson:"author_name"`
	Text       string    `json:"text" bson:"text"`
	Creation   time.Time `json:"creation" bson:"creation,omitempty"`
}

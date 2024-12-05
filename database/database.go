package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type DatabaseInterface interface {
	CreateUser(user interface{}) error
	FindUserByUsername(username string, user interface{}) error
}

type MongoDatabase struct {
	DB *mongo.Database
}

func (m *MongoDatabase) CreateUser(user interface{}) error {
	collection := m.DB.Collection("users")
	_, err := collection.InsertOne(context.Background(), user)
	return err
}

func (m *MongoDatabase) FindUserByUsername(username string, user interface{}) error {
	collection := m.DB.Collection("users")
	err := collection.FindOne(context.Background(), map[string]string{"username": username}).Decode(user)
	return err
}

type GormDatabase struct {
	DB *gorm.DB
}

func (g *GormDatabase) CreateUser(user interface{}) error {
	return g.DB.Create(user).Error
}

func (g *GormDatabase) FindUserByUsername(username string, user interface{}) error {
	return g.DB.Where("username = ?", username).First(user).Error
}

func NewDatabaseService(db interface{}) DatabaseInterface {
	switch v := db.(type) {
	case *mongo.Database:
		return &MongoDatabase{DB: v}
	case *gorm.DB:
		return &GormDatabase{DB: v}
	default:
		panic("unsupported database type")
	}
}

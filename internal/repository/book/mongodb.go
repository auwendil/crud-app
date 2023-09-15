package book

import (
	"context"
	"github.com/auwendil/crud-app/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoDBRepo struct {
	collection *mongo.Collection
}

const (
	mongoDBName         = "db"
	mongoCollectionName = "books"
)

func NewMongoDBRepo(connectionURI string) (*MongoDBRepo, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionURI))
	if err != nil {
		return nil, err
	}

	mongoDB := &MongoDBRepo{
		collection: client.Database(mongoDBName).Collection(mongoCollectionName),
	}

	return mongoDB, nil
}

func (r *MongoDBRepo) GetAllBooks() ([]*models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	allValuesFilter := bson.D{}
	cursor, err := r.collection.Find(ctx, allValuesFilter)
	if err != nil {
		return nil, err
	}

	ctx, cancelCursor := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCursor()

	books := []*models.Book{}
	if err = cursor.All(ctx, &books); err != nil {
		return nil, err
	}

	return books, nil
}

func (r *MongoDBRepo) GetBook(id string) (*models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{"_id", objID}}
	result := r.collection.FindOne(ctx, filter)
	if err = result.Err(); err != nil {
		return nil, err
	}

	var book *models.Book
	if err = result.Decode(&book); err != nil {
		return nil, err
	}

	return book, nil
}

func (r *MongoDBRepo) AddBook(b *models.Book) (*models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()

	bytes, err := bson.Marshal(b)
	if err != nil {
		return nil, err
	}

	result, err := r.collection.InsertOne(ctx, bytes)
	if err != nil {
		return nil, err
	}

	var id string
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		id = oid.Hex()
	}

	b.ID = id
	return b, nil
}

func (r *MongoDBRepo) UpdateBook(id string, updatedBook *models.Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	updatedBook.UpdatedAt = time.Now()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", objID}}
	updatedObject := bson.M{"$set": updatedBook}

	res := r.collection.FindOneAndUpdate(ctx, filter, updatedObject)
	if err = res.Err(); err != nil {
		return err
	}

	return nil
}

func (r *MongoDBRepo) DeleteBook(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", objID}}
	_, err = r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoDBRepo) DeleteAllBooks() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	takeAllFilter := bson.D{}
	_, err := r.collection.DeleteMany(ctx, takeAllFilter)
	if err != nil {
		return err
	}

	return nil
}

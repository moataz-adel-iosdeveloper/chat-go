package userRepository

import (
	"chat-go/internal/config"
	"chat-go/internal/models"
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// helper to get the users collection at runtime. Returns an error if DB not initialized.
func getUserCollection() (*mongo.Collection, error) {
	col := config.GetUserCollection()
	if col == nil {
		return nil, errors.New("database not initialized: call services.InitDB before using repositories")
	}
	return col, nil
}

func CreateUser(user *models.User) (*models.User, error) {
	col, err := getUserCollection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = col.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateToken(user *models.User) error {
	col, err := getUserCollection()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": user.ID}

	update := bson.M{
		"$set": bson.M{
			"token": user.Token,
		},
	}

	result, err := col.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("Error updating user:", err)
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func FindUserByEmail(email string) (*models.User, error) {
	col, err := getUserCollection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user *models.User
	err = col.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, nil
	}
	return user, err
}

func FindUserByID(id string) (*models.User, error) {
	col, err := getUserCollection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user *models.User
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = col.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetAllUsers() ([]*models.User, error) {
	col, err := getUserCollection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var users []*models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

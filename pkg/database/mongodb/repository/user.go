package repository

import (
	"context"
	"errors"

	"github.com/organization_api/pkg/database"
	"github.com/organization_api/pkg/database/mongodb/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository represents the MongoDB collection for user data.
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository initializes a new UserRepository instance.
func NewUserRepository() *UserRepository {
	// Get the MongoDB collection for user data.
	db := database.GetDatabase()
	return &UserRepository{collection: db.Collection("user")}
}

// FindUserByEmail retrieves a user from the database by their email address.
func (repo *UserRepository) FindUserByEmail(email string) (*models.User, error) {
	// Search for the user by email.
	filter := bson.M{"email": email}
	var user models.User
	err := repo.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		// Return nil if no user is found, otherwise, return an error.
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// CreateUser inserts a new user into the database.
func (repo *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	// Check if the user already exists by email.
	existingUser := &models.User{}
	err := repo.collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(existingUser)
	if err == nil {
		// Return an error if the email already exists.
		return nil, errors.New("email already exists")
	}

	// Insert the new user into the database.
	createdUser, err := repo.collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": createdUser.InsertedID}
	var insertedUser models.User
	err = repo.collection.FindOne(context.TODO(), filter).Decode(&insertedUser)
	if err != nil {
		return nil, err
	}

	return &insertedUser, nil
}

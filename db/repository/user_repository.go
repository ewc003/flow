package repository

import (
	"context"
	"errors"
	"server/db/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	BaseRepository
}

func NewUserRepository(db *mongo.Client) *UserRepository {
	return &UserRepository{
		BaseRepository{
			Collection: db.Database("flow").Collection("users"),
		},
	}
}

// CreateUser creates a new user with hashed password
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	// Check if email already exists
	existingUser := &models.User{}
	err := r.FindOne(ctx, bson.M{"email": user.Email}, existingUser)
	if err == nil {
		return errors.New("user with this email already exists")
	}

	// Check if username already exists
	err = r.FindOne(ctx, bson.M{"username": user.Username}, existingUser)
	if err == nil {
		return errors.New("username already taken")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Set timestamps
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Update the password with the hashed version
	user.Password = string(hashedPassword)

	// Create the user
	return r.Create(ctx, user)
}

// FindUserByEmail finds a user by email address
func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := r.FindOne(ctx, bson.M{"email": email}, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindUserByID finds a user by ID
func (r *UserRepository) FindUserByID(ctx context.Context, id string) (*models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	user := &models.User{}
	err = r.FindOne(ctx, bson.M{"_id": objectID}, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ValidateCredentials checks if the provided email and password match a user
func (r *UserRepository) ValidateCredentials(ctx context.Context, email, password string) (*models.User, error) {
	user, err := r.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

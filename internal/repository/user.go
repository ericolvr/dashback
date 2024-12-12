package repository

import (
	"context"

	"github.com/Alarmtekgit/websocket/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	FindUserByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error)
	FindUserByMobile(ctx context.Context, mobile string) (*domain.User, error)
	DeleteUser(ctx context.Context, id primitive.ObjectID) error
	GetAllUsers(ctx context.Context) ([]domain.User, error)
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{
		collection: db.Collection("users"),
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	hashedPassword, err := domain.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	_, err = r.collection.InsertOne(ctx, user)
	return err
}

func (r *userRepository) FindUserByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	return &user, err
}

func (r *userRepository) FindUserByMobile(ctx context.Context, mobile string) (*domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"mobile": mobile}).Decode(&user)
	return &user, err
}

func (r *userRepository) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

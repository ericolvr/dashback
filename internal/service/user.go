package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Alarmtekgit/websocket/internal/domain"
	"github.com/Alarmtekgit/websocket/internal/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, message *domain.User) error
	FindUserByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error)
	FindUserByMobile(ctx context.Context, mobile string) (*domain.User, error)
	DeleteUser(ctx context.Context, id primitive.ObjectID) error
	GetAllUsers(ctx context.Context) ([]domain.User, error)
	Authenticate(ctx context.Context, mobile, password string) (*domain.AuthResponse, error)
}

var jwtSecret = []byte(viper.GetString("JWT-SECRET"))

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, user *domain.User) error {
	user.Timestamp = time.Now()
	return s.repo.CreateUser(ctx, user)
}

func (s *userService) FindUserByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	return s.repo.FindUserByID(ctx, id)
}

func (s *userService) FindUserByMobile(ctx context.Context, mobile string) (*domain.User, error) {
	return s.repo.FindUserByMobile(ctx, mobile)
}

func (s *userService) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *userService) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	return s.repo.GetAllUsers(ctx)
}

func (s *userService) Authenticate(ctx context.Context, mobile, password string) (*domain.AuthResponse, error) {
	user, err := s.repo.FindUserByMobile(ctx, mobile)
	if err != nil {
		return nil, errors.New("invalid mobile or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid mobile or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name":   user.Name,
		"mobile": user.Mobile,
		"role":   user.Role,
	})
	fmt.Printf("user data: %v\n", token)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Name:  user.Name,
		Token: tokenString,
		Role:  user.Role,
	}, nil

}

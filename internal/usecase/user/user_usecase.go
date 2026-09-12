package user

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	entity "be_evindo/internal/entity/user"
)

type UserUsecase struct {
	repository entity.Repository
	jwtSecret  []byte
	tokenTTL   time.Duration
}

func NewUserUsecase(repository entity.Repository, jwtSecret string) *UserUsecase {
	return &UserUsecase{
		repository: repository,
		jwtSecret:  []byte(jwtSecret),
		tokenTTL:   24 * time.Hour,
	}
}

func (usecase *UserUsecase) Register(ctx context.Context, name, email, password string) (*entity.User, error) {
	return usecase.CreateUser(ctx, name, email, password, "USER")
}

func (usecase *UserUsecase) CreateUser(ctx context.Context, name, email, password, role string) (*entity.User, error) {
	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))
	role = strings.ToUpper(strings.TrimSpace(role))
	if name == "" || email == "" || password == "" {
		return nil, errors.New("name, email, and password are required")
	}
	if role != "USER" && role != "APPROVER" {
		return nil, errors.New("role must be USER or APPROVER")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	newUser := &entity.User{Name: name, Email: email, Password: string(hashedPassword), Role: role}
	if err := usecase.repository.Create(ctx, newUser); err != nil {
		return nil, err
	}
	newUser.Password = ""
	return newUser, nil
}

func (usecase *UserUsecase) Login(ctx context.Context, email, password string) (string, *entity.User, error) {
	account, err := usecase.repository.FindByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
	if err != nil || bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password)) != nil {
		return "", nil, errors.New("invalid email or password")
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"user_id": account.ID,
		"role":    account.Role,
		"iat":     now.Unix(),
		"exp":     now.Add(usecase.tokenTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(usecase.jwtSecret)
	if err != nil {
		return "", nil, err
	}
	account.Password = ""
	return signedToken, account, nil
}

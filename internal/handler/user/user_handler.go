package user

import (
	"errors"
	"net/http"

	usecase "be_evindo/internal/usecase/user"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserHandler struct {
	usecase *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: userUsecase}
}

func (handler *UserHandler) Register(context *gin.Context) {
	var request RegisterRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account, err := handler.usecase.CreateUser(context.Request.Context(), request.Name, request.Email, request.Password, request.Role)
	if err != nil {
		var postgresError *pgconn.PgError
		if errors.As(err, &postgresError) && postgresError.Code == "23505" {
			context.JSON(http.StatusConflict, gin.H{"error": "email is already registered"})
			return
		}
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusCreated, UserResponse{
		ID: account.ID, Name: account.Name, Email: account.Email, Role: account.Role,
	})
}

func (handler *UserHandler) Login(context *gin.Context) {
	var request LoginRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, account, err := handler.usecase.Login(context.Request.Context(), request.Email, request.Password)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User: UserResponse{
			ID: account.ID, Name: account.Name, Email: account.Email, Role: account.Role,
		},
	})
}

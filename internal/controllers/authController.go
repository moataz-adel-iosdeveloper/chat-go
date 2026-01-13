package controllers

import (
	"chat-go/internal/helpers"
	"chat-go/internal/models"
	userRepository "chat-go/internal/repositories/user"

	"log"

	"chat-go/internal/services"
	"io"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Login handles authentication logic and returns a JSON-serializable
// response, an HTTP status code, and an error if one occurred.
// This is a minimal placeholder implementation.
func Login(req *http.Request) (response models.APIResponse, statusCode int) {
	var lr models.UserCredentials
	if err := helpers.QueryParamsToStruct(req, &lr); err != nil {
		log.Println("Error parsing query/form params:", err)
		return models.ErrorResponse("invalid query/form params", nil), http.StatusBadRequest
	}

	if lr.Email == "" || lr.Password == "" {
		return models.ErrorResponse("missing username or password", nil), http.StatusBadRequest
	}

	// check E-mail is existing
	user, err := userRepository.FindUserByEmail(lr.Email)
	if err != nil || user == nil {
		return models.ErrorResponse("invalid email or password", nil), http.StatusUnauthorized
	}

	// get password papper
	papperString := os.Getenv("PASSWORD_PEPPER")
	if papperString == "" {
		return models.ErrorResponse("server misconfiguration", nil), http.StatusInternalServerError
	}

	// check password
	err = services.ComparePassword(user.Password, lr.Password, papperString)
	if err != nil {
		return models.ErrorResponse("invalid email or password", nil), http.StatusUnauthorized
	}

	// create JWT token
	token := services.GenerateToken(user.ID.Hex())
	if token == nil {
		return models.ErrorResponse("failed to generate token", nil), http.StatusInternalServerError
	}

	user.Token = string(token)

	// update user with token
	err = userRepository.UpdateToken(user)
	if err != nil {
		return models.ErrorResponse("failed to update user token", nil), http.StatusInternalServerError
	}

	userResp := models.UserResponse{
		ID:       user.ID.Hex(),
		Username: user.Username,
		Email:    user.Email,
		Token:    user.Token,
	}

	resp := map[string]models.UserResponse{
		"user": userResp,
	}

	return models.SuccessResponse("login successful", resp), http.StatusOK
}

func Register(req *http.Request) (response models.APIResponse, statusCode int) {
	var ur models.UserRegistration

	if err := helpers.QueryParamsToStruct(req, &ur); err != nil {
		if err == io.EOF {
			return models.ErrorResponse("empty body", nil), http.StatusBadRequest
		}
		log.Println("Error parsing query params:", err)
		return models.ErrorResponse("invalid query params", nil), http.StatusBadRequest
	}

	if ur.Username == "" || ur.Email == "" || ur.Password == "" {
		return models.ErrorResponse("missing username, email or password", nil), http.StatusBadRequest
	}

	// check E-mail is existing
	existingUser, err := userRepository.FindUserByEmail(ur.Email)

	if err != nil {
		log.Println("Error checking existing user:", err)
		return models.ErrorResponse("server error", nil), http.StatusInternalServerError
	}
	if existingUser != nil {
		log.Println("Error checking existing user:")
		return models.ErrorResponse("email already in use", nil), http.StatusConflict
	}

	// validate email format
	if !helpers.CheckEmailFormat(ur.Email) {
		return models.ErrorResponse("invalid email format", nil), http.StatusBadRequest
	}

	// get password papper
	papperString := os.Getenv("PASSWORD_PEPPER")
	if papperString == "" {
		return models.ErrorResponse("server misconfiguration", nil), http.StatusInternalServerError
	}
	hashedPassword, err := services.HashPassword(ur.Password, papperString)
	if err != nil {
		return models.ErrorResponse("failed to hash password", nil), http.StatusInternalServerError
	}
	// generate new user ID
	var id = primitive.NewObjectID()

	// create token
	token := services.GenerateToken(id.Hex())
	if token == nil {
		return models.ErrorResponse("failed to generate token", nil), http.StatusInternalServerError
	}

	newUser := &models.User{
		ID:       id,
		Username: ur.Username,
		Email:    ur.Email,
		Password: hashedPassword,
		Token:    string(token),
	}

	// create user
	createdUser, err := userRepository.CreateUser(newUser)
	if err != nil {
		return models.ErrorResponse("failed to create user", nil), http.StatusInternalServerError
	}

	userResp := models.UserResponse{
		ID:       createdUser.ID.Hex(),
		Username: createdUser.Username,
		Email:    createdUser.Email,
		Token:    createdUser.Token,
	}

	resp := map[string]models.UserResponse{
		"user": userResp,
	}
	return models.SuccessResponse("registration successful", resp), http.StatusCreated
}

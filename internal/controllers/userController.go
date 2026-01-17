package controllers

import (
	"chat-go/internal/models"
	userRepository "chat-go/internal/repositories/user"
	"log"
	"net/http"
)

func GetAllUsers(req *http.Request) (response models.APIResponse, statusCode int) {

	// Get the authenticated user ID from context
	authUserID := req.Context().Value("userID")
	if authUserID == nil {
		log.Println("Error parsing auth ID")
		return models.ErrorResponse("Error parsing auth ID", nil), http.StatusBadRequest
	}
	var allUsers []*models.User
	allUsers, err := userRepository.GetAllUsers()
	if err != nil {
		return models.ErrorResponse("failed to get users", nil), http.StatusInternalServerError
	}
	var UserResponses []models.UserResponse
	for _, user := range allUsers {
		if user.ID.Hex() == authUserID {
			continue
		}
		UserResponses = append(UserResponses, models.UserResponse{
			ID:       user.ID.Hex(),
			Username: user.Username,
			Email:    user.Email,
		})
	}
	return models.SuccessResponse("users retrieved successfully", UserResponses), http.StatusOK
}

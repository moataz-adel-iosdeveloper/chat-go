package controllers

import (
	"chat-go/internal/models"
	userRepository "chat-go/internal/repositories/user"
	"net/http"
)

func GetAllUsers(req *http.Request) (response models.APIResponse, statusCode int) {
	var allUsers []*models.User
	allUsers, err := userRepository.GetAllUsers()
	if err != nil {
		return models.ErrorResponse("failed to get users", nil), http.StatusInternalServerError
	}
	var UserResponses []models.UserResponse
	for _, user := range allUsers {
		UserResponses = append(UserResponses, models.UserResponse{
			ID:       user.ID.Hex(),
			Username: user.Username,
			Email:    user.Email,
		})
	}
	return models.SuccessResponse("users retrieved successfully", UserResponses), http.StatusOK
}

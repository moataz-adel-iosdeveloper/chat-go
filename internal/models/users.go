package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string             `json:"username"`
	Email    string             `json:"email"`
	Password string             `json:"-" bson:"password"`
	Token    string             `json:"token,omitempty" bson:"token,omitempty"`
}

type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegistration struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token,omitempty"`
}

package handler

import "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"

type RegisterResponse struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
	Token   string `json:"token"`
}

type LoginResponse struct {
	Message string       `json:"message"`
	Token   string       `json:"token"`
	User    *domain.User `json:"user"`
}

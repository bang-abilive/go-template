package user

import "time"

type RegisterUserResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email"`
	ID        int64     `json:"id"`
}

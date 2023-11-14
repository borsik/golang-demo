package user

import (
	"github.com/google/uuid"
	"time"
)

// User holds the structure of user entity in db, all fields are shown in response json except password
type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Nickname  string    `json:"nickname"`
	Password  string    `json:"-"`
	Email     string    `json:"email"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// InputUser represents json body for users POST/PUT api
type InputUser struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Nickname  string `json:"nickname" validate:"required"`
	Password  string `json:"password" validate:"required,ascii,min=8,max=72"`
	Email     string `json:"email" validate:"required,email"`
	Country   string `json:"country" validate:"required,iso3166_1_alpha2"`
}

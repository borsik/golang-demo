package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

func (u *User) BeforeCreate(_ *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Nickname  string    `json:"nickname" gorm:"index:idx_nickname,unique"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

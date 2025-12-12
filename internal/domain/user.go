package domain

import "movie-app/internal/enums"

type User struct {
	BaseModel
	Name     string     `gorm:"type:varchar(100);not null" json:"name"`
	Email    string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password string     `gorm:"type:varchar(255);not null" json:"-"`
	Role     enums.Role `gorm:"type:varchar(20);default:'user'" json:"role"` // Menggunakan Enum
}

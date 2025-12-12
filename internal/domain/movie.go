package domain

type Movie struct {
	BaseModel
	Title       string `gorm:"type:varchar(255);not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	Duration    int    `gorm:"not null" json:"duration"`
	Genre       string `gorm:"type:varchar(100)" json:"genre"`
	PosterURL   string `gorm:"type:varchar(255)" json:"poster_url"`
}

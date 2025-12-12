package domain

type Studio struct {
	BaseModel
	Name     string `gorm:"type:varchar(100);not null" json:"name"`
	Capacity int    `gorm:"not null" json:"capacity"`
	Seats    []Seat `gorm:"foreignKey:StudioID" json:"seats,omitempty"`
}

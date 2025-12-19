package model

import (
	"time"
)

type NutritionGoal struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID        string    `gorm:"type:uuid;index;not null" json:"user_id"`
	Calories      float64   `json:"calories"`
	Protein       float64   `json:"protein"`
	Carbohydrates float64   `json:"carbohydrates"`
	Fat           float64   `json:"fat"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

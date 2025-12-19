package model

import (
	"time"
)

type Food struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name          string    `gorm:"index;not null" json:"name"`
	Calories      float64   `json:"calories"`
	Protein       float64   `json:"protein"`
	Carbohydrates float64   `json:"carbohydrates"`
	Fat           float64   `json:"fat"`
	CreatedAt     time.Time `json:"created_at"`
}

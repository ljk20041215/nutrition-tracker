package model

import (
	"time"
)

// FoodRecord 食物记录模型
type FoodRecord struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	MealRecordID  string    `gorm:"type:uuid;index;not null" json:"meal_record_id"` // 关联的餐次ID
	FoodID        string    `gorm:"type:uuid;index;not null" json:"food_id"`        // 关联的食物ID
	FoodName      string    `gorm:"type:varchar(100);not null" json:"food_name"`    // 冗余存储食物名称，提高查询效率
	Quantity      float64   `gorm:"type:float;not null" json:"quantity"`            // 份量
	Unit          string    `gorm:"type:varchar(20);not null" json:"unit"`          // 单位（g, kg, ml, 个等）
	Calories      float64   `json:"calories"`                                        // 实际摄入的热量（根据份量计算）
	Protein       float64   `json:"protein"`                                         // 实际摄入的蛋白质
	Carbohydrates float64   `json:"carbohydrates"`                                   // 实际摄入的碳水化合物
	Fat           float64   `json:"fat"`                                             // 实际摄入的脂肪
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// 关联关系
	MealRecord MealRecord `gorm:"foreignKey:MealRecordID" json:"-"`
	Food       Food       `gorm:"foreignKey:FoodID" json:"-"`
}

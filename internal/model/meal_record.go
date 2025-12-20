package model

import (
	"encoding/json"
	"errors"
	"time"
)

// MealType 餐次类型常量
type MealType int

const (
	Breakfast MealType = iota + 1 // 早餐
	Lunch                         // 午餐
	Dinner                        // 晚餐
	Snack                         // 加餐
)

// MealTypeMap 餐次类型映射
type MealTypeMap struct {
	IntValue    MealType
	StringValue string
}

// MealTypeValues 餐次类型值列表
var MealTypeValues = map[string]MealType{
	"breakfast": Breakfast,
	"lunch":     Lunch,
	"dinner":    Dinner,
	"snack":     Snack,
}

// MealTypeStrings 餐次类型字符串映射
var MealTypeStrings = map[MealType]string{
	Breakfast: "breakfast",
	Lunch:     "lunch",
	Dinner:    "dinner",
	Snack:     "snack",
}

// MarshalJSON 实现MealType的JSON序列化方法
func (m MealType) MarshalJSON() ([]byte, error) {
	if str, exists := MealTypeStrings[m]; exists {
		return json.Marshal(str)
	}
	return nil, errors.New("无效的餐次类型")
}

// UnmarshalJSON 实现MealType的JSON反序列化方法
func (m *MealType) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		// 如果不是字符串，尝试解析为数字
		var num int
		if err := json.Unmarshal(data, &num); err != nil {
			return err
		}
		*(*int)(m) = num
		return nil
	}

	if val, exists := MealTypeValues[str]; exists {
		*m = val
		return nil
	}

	return errors.New("无效的餐次类型字符串")
}

// MealRecord 餐次记录模型
type MealRecord struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    string    `gorm:"type:uuid;index;not null" json:"user_id"`
	Date      time.Time `gorm:"type:date;index;not null" json:"date"`
	MealType  MealType  `gorm:"type:int;not null" json:"meal_type"` // 1:早餐, 2:午餐, 3:晚餐, 4:加餐
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}